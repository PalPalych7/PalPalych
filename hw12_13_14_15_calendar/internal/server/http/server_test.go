package internalhttp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	st "github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

// перед тестами запустить прилжение, чтобы стартовался сервер
func TestHTTP(t *testing.T) {
	t.Run("main", func(t *testing.T) {
		client := http.Client{
			Timeout: time.Second * 5,
		}
		// Создание заявки
		reqBody := ForCreate{}
		reqBody.Details = "newDetails"
		reqBody.StartDate = "21.01.2023"
		reqBody.Title = "newTitle"
		bodyRaw, err := json.Marshal(reqBody)
		req, err := http.NewRequest(http.MethodPost, "http://127.0.0.102:5000/CreateEvent", bytes.NewBuffer(bodyRaw))
		require.NoError(t, err)
		_, errResp := client.Do(req)
		require.NoError(t, errResp)

		// проверка колическв заявок по дате
		reqBody2 := StartDate{}
		reqBody2.StartDateStr = "21.01.2023"
		bodyRaw2, err2 := json.Marshal(reqBody2)
		req2, err2 := http.NewRequest(http.MethodPost, "http://127.0.0.102:5000/GetEventByDate", bytes.NewBuffer(bodyRaw2))
		require.NoError(t, err2)
		resp2, errResp2 := client.Do(req2)
		require.NoError(t, errResp2)
		bodyBytes, errrr := ioutil.ReadAll(resp2.Body)
		require.NoError(t, errrr)
		var eventList []st.Event
		errUnm := json.Unmarshal(bodyBytes, &eventList)
		require.NoError(t, errUnm)
		myId := eventList[0].ID
		require.Equal(t, 1, len(eventList))

		// обновление заявки
		reqBody3 := ForUpdate{}
		reqBody3.EventID = myId
		reqBody3.Details = "Details2"
		reqBody3.StartDate = "21.01.2023"
		reqBody3.Title = "newTitle"
		bodyRaw3, err3 := json.Marshal(reqBody2)
		req3, err3 := http.NewRequest(http.MethodPost, "http://127.0.0.102:5000/UpdateEvent", bytes.NewBuffer(bodyRaw3))
		require.NoError(t, err3)
		_, errResp3 := client.Do(req3)
		require.NoError(t, errResp3)

		// удаление заявки
		reqBody4 := ForDelete{}
		reqBody4.EventID = myId
		bodyRaw4, err4 := json.Marshal(reqBody4)
		req4, err4 := http.NewRequest(http.MethodPost, "http://127.0.0.102:5000/DeleteEvent", bytes.NewBuffer(bodyRaw4))
		require.NoError(t, err4)
		_, errResp4 := client.Do(req4)
		require.NoError(t, errResp4)

		// повторная проверка колическв заявок по дате
		resp2, errResp2 = client.Do(req2)
		require.NoError(t, errResp2)
		bodyBytes, errrr = ioutil.ReadAll(resp2.Body)
		require.NoError(t, errrr)
		errUnm = json.Unmarshal(bodyBytes, &eventList)
		require.NoError(t, errUnm)
		require.Equal(t, 0, len(eventList))

	})
}
