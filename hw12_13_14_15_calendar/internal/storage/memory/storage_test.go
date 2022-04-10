package memorystorage

import (
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

		// update event
		err = storage.UpdateEvent("123", "t88", "15.3.2029", "xxxfdsfsd", 2)
		require.Error(t, err)
		require.ErrorIs(t, err, IdNotFound)

		var myKey string
		for Key := range storage.Events {
			myKey = Key
			break
		}

		err = storage.UpdateEvent(myKey, "t999", "15.3.2029", "xxxfdsfsd", 2)
		require.NoError(t, err)
		require.Equal(t, "t999", storage.Events[myKey].Title)

		// Get Event By Date
		myEventList, err2 := storage.GetEventByDate("15.3.2028") // not in calendar
		require.NoError(t, err2)
		require.Equal(t, 0, len(myEventList))

		myEventList, err2 = storage.GetEventByDate("15.3.2029") // in calendar
		require.NoError(t, err2)
		require.Equal(t, 1, len(myEventList))

		// Get Event By Month
		myEventList, err2 = storage.GetEventMonth("15.3.2028") // not begin month
		require.Error(t, err2)
		require.ErrorIs(t, err2, ErrNotBeginMonth)

		myEventList, err2 = storage.GetEventMonth("1.3.2028") // not in calendar
		require.NoError(t, err2)
		require.Equal(t, 0, len(myEventList))

		myEventList, err2 = storage.GetEventMonth("1.3.2022") // 2 event in calendar
		require.NoError(t, err2)
		require.Equal(t, 2, len(myEventList))

		// Get Event By Week
		myEventList, err2 = storage.GetEventWeek("29.3.2022") // not begin month
		require.Error(t, err2)
		require.ErrorIs(t, err2, ErrNotBeginWeek)

		myEventList, err2 = storage.GetEventWeek("28.3.2022") // not in calendar
		require.NoError(t, err2)
		require.Equal(t, 0, len(myEventList))

		myEventList, err2 = storage.GetEventWeek("14.3.2022") // 1 event in calendar
		require.NoError(t, err2)
		require.Equal(t, 1, len(myEventList))

		// delete event
		err = storage.DeleteEvent("123")
		require.Error(t, err)
		require.ErrorIs(t, err, IdNotFound)

		err = storage.DeleteEvent(myKey)
		require.NoError(t, err)
		require.Equal(t, 2, len(storage.Events))

	})
}
