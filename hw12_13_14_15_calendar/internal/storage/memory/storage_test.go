package memorystorage

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("main", func(t *testing.T) {
		storage := New()

		// new event
		err := storage.CreateEvent("loko", "8.1.2028", "lokomotiv", 1)
		require.NoError(t, err)

		err = storage.CreateEvent("t2", "2.3.2022", "xxx", 2)
		require.NoError(t, err)

		err = storage.CreateEvent("t3", "15.3.2022", "xxxfdsfsd", 2)
		require.NoError(t, err)

		err = storage.CreateEvent("t3", "33.3.2011", "xxxfdsf", 1)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrDate)

		require.Equal(t, 3, len(storage.Events))

		// Get Event By Date
		myEventList, err2 := storage.GetEventByDate("15.3.2021") // not in calendar
		require.NoError(t, err2)
		require.Equal(t, 0, len(myEventList))

		// update event
		err10 := storage.UpdateEvent("123", "t88", "15.3.2029", "xxxfdsfsd", 2)
		require.Error(t, err10)
		require.ErrorIs(t, err10, EventIDIsNotFound)

		myEventList3, err3 := storage.GetEventByDate("15.3.2022") // in calendar
		require.NoError(t, err3)
		require.Equal(t, 1, len(myEventList3))

		// Get Event By Month
		_, err4 := storage.GetEventMonth("15.3.2028") // not begin month
		require.Error(t, err4)
		require.ErrorIs(t, err4, ErrNotBeginMonth)

		myEventList5, err5 := storage.GetEventMonth("1.3.2028") // not in calendar
		require.NoError(t, err5)
		require.Equal(t, 0, len(myEventList5))

		myEventList6, err6 := storage.GetEventMonth("1.3.2022") // 2 event in calendar
		fmt.Println(myEventList6)
		require.NoError(t, err6)
		require.Equal(t, 2, len(myEventList6))

		// Get Event By Week
		_, err7 := storage.GetEventWeek("29.3.2022") // not begin month
		require.Error(t, err7)
		require.ErrorIs(t, err7, ErrNotBeginWeek)

		myEventList8, err8 := storage.GetEventWeek("28.3.2022") // not in calendar
		require.NoError(t, err8)
		require.Equal(t, 0, len(myEventList8))

		myEventList9, err9 := storage.GetEventWeek("14.3.2022") // 1 event in calendar
		require.NoError(t, err9)
		require.Equal(t, 1, len(myEventList9))

		// delete event
		err11 := storage.DeleteEvent("123")
		require.Error(t, err11)
		require.ErrorIs(t, err11, EventIDIsNotFound)

		var myKey string
		for Key := range storage.Events {
			myKey = Key
			break
		}
		err12 := storage.UpdateEvent(myKey, "t999", "15.3.2029", "xxxfdsfsd", 2)
		require.NoError(t, err12)
		require.Equal(t, "t999", storage.Events[myKey].Title)

		err13 := storage.DeleteEvent(myKey)
		require.NoError(t, err13)
		require.Equal(t, 2, len(storage.Events))
	})
}
