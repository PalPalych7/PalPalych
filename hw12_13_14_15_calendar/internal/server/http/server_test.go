package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	st "github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

//  Перед тестами запустить прилжение, чтобы стартовался сервер
func TestHTTP(t *testing.T) {
	t.Run("main", func(t *testing.T) {
		client := http.Client{
			Timeout: time.Second * 5,
		}
		constDate := "21.01.2023"
		myHTTP := "http://127.0.0.103:5000/"
		ctx := context.Background()
		// Создание заявки
		reqBody := ForCreate{}
		reqBody.Details = "newDetails"
		reqBody.StartDate = constDate
		reqBody.Title = "newTitle"
		reqBody.UserID = 123
		bodyRaw, errM := json.Marshal(reqBody)
		require.NoError(t, errM)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, myHTTP+"CreateEvent", bytes.NewBuffer(bodyRaw))
		require.NoError(t, err)
		_, errResp := client.Do(req) //nolint
		require.NoError(t, errResp)

		// Проверка колическв заявок по дате
		reqBody2 := StartDate{}
		reqBody2.StartDateStr = constDate
		bodyRaw2, errM2 := json.Marshal(reqBody2)
		require.NoError(t, errM2)
		req2, err2 := http.NewRequestWithContext(ctx, http.MethodPost, myHTTP+"GetEventByDate", bytes.NewBuffer(bodyRaw2))
		require.NoError(t, err2)
		resp2, errResp2 := client.Do(req2)
		require.NoError(t, errResp2)
		defer resp2.Body.Close()
		bodyBytes, errrr := ioutil.ReadAll(resp2.Body)
		require.NoError(t, errrr)
		var eventList []st.Event
		errUnm := json.Unmarshal(bodyBytes, &eventList)
		require.NoError(t, errUnm)
		myID := eventList[0].ID
		require.Equal(t, 1, len(eventList))

		// обновление заявки
		reqBody3 := ForUpdate{}
		reqBody3.EventID = myID
		reqBody3.Details = "Details2"
		reqBody3.StartDate = constDate
		reqBody3.Title = "newTitle"
		bodyRaw3, errM3 := json.Marshal(reqBody3)
		require.NoError(t, errM3)

		req3, err3 := http.NewRequestWithContext(ctx, http.MethodPost, myHTTP+"UpdateEvent", bytes.NewBuffer(bodyRaw3))
		require.NoError(t, err3)
		_, errResp3 := client.Do(req3) //nolint
		require.NoError(t, errResp3)

		// удаление заявки
		reqBody4 := ForDelete{}
		reqBody4.EventID = myID
		bodyRaw4, errM4 := json.Marshal(reqBody4)
		require.NoError(t, errM4)

		req4, err4 := http.NewRequestWithContext(ctx, http.MethodPost, myHTTP+"DeleteEvent", bytes.NewBuffer(bodyRaw4))
		require.NoError(t, err4)
		_, errResp4 := client.Do(req4) //nolint
		require.NoError(t, errResp4)

		// повторная проверка колическв заявок по дате
		resp2, errResp2 = client.Do(req2)
		require.NoError(t, errResp2)
		defer resp2.Body.Close()
		bodyBytes, errrr = ioutil.ReadAll(resp2.Body)
		require.NoError(t, errrr)
		errUnm = json.Unmarshal(bodyBytes, &eventList)
		require.NoError(t, errUnm)
		require.Equal(t, 0, len(eventList))
	})
}
