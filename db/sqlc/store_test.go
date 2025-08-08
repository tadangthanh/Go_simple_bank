package sqlc

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	store := NewStore(conn)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before: ", account1.Balance, account2.Balance)
	// run n concurrent transfer transactions
	n := 5
	errs := make(chan error)
	results := make(chan TransferTxResult)
	amount := int64(10)
	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}
	// check results
	existed := make(map[int]bool) //existed là map dùng để track các giá trị k đã xuất hiện.Đảm bảo mỗi giá trị k chỉ xảy ra một lần.
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)
		ToEntry := result.ToEntry
		require.NotEmpty(t, ToEntry)
		require.Equal(t, account2.ID, ToEntry.AccountID)
		require.Equal(t, amount, ToEntry.Amount)
		require.NotZero(t, ToEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), ToEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		// check balance
		fmt.Println(">> tx: ", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance //số tiền giảm của tài khoản gửi.
		diff2 := toAccount.Balance - account2.Balance   // số tiền tăng của tài khoản nhận.
		require.Equal(t, diff1, diff2)                  //tiền mất ở tài khoản gửi == tiền nhận ở tài khoản đích
		require.True(t, diff1 > 0)                      //Kiểm tra rằng có thay đổi số dư thực sự xảy ra (phải > 0), tức là transaction đã chạy.
		require.True(t, diff1%amount == 0)              //đảm bảo rằng tổng số tiền bị trừ là bội số của amount – tức là kết quả của k lần transaction.
		k := int(diff1 / amount)                        //Tính ra số lần transaction đã chạy thành công cho account này, bằng cách chia số tiền thay đổi cho amount.
		require.True(t, k >= 1 && k <= n)               //Xác thực rằng k nằm trong phạm vi [1, n] – không vượt quá số lần test concurrent.
		require.NotContains(t, existed, k)
		existed[k] = true //Điều này kiểm tra rằng mỗi transaction là duy nhất và không có transaction bị lặp lại hoặc chạy trùng.
	}
	// check the final updated balances
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println(">> after: ", updateAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)
}
