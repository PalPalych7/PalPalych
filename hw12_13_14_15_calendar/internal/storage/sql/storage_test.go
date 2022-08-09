package sqlstorage

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("main", func(t *testing.T) {
		// Тесты будут не столь подробные, по сравнению с memory
		// важно проверить факт успешного конекта-дисконекта и работу основной логики
		storage := New("hw12", "otus1", "123456")

		// DB connect
		err := storage.Connect(context.Background())
		require.NoError(t, err)
		err = storage.DBConnect.Ping()
		require.NoError(t, err)

		// Get Event By Date
		//		myEventList, err2 := storage.GetEventByDate("2022-05-11")
		myEventList, err2 := storage.GetEventByDate("11.05.2022")
		require.NoError(t, err2)
		len1 := len(myEventList)

		// new event
		err = storage.CreateEvent("t1", "11.05.2022", "something", 1)
		require.NoError(t, err)

		// Get Event By Date2
		myEventList, err2 = storage.GetEventByDate("11.05.2022")
		require.NoError(t, err2)
		len2 := len(myEventList)
		fmt.Println(len2)
		require.Equal(t, len1+1, len2)

		// close DB
		err = storage.Close(context.Background())
		require.NoError(t, err)
		err = storage.DBConnect.Ping()
		require.Error(t, err)
	})
}
