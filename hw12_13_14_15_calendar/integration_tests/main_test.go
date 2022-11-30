package integration_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	_ "github.com/jackc/pgx/stdlib" // justifying
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type mySuite struct {
	suite.Suite
	ctx       context.Context
	client    http.Client
	hostName  string
	DBConnect *sql.DB
}

type StartDate struct {
	StartDateStr string
}

type ForCreate struct {
	Title     string
	StartDate string
	Details   string
	UserID    uint32
}

type ForUpdate struct {
	EventID   string
	Title     string
	StartDate string
	Details   string
	UserID    uint32
}

type ForDelete struct {
	EventID string
}

type OneEvent struct {
	EventID   string
	Title     string
	StartDate string
	Details   string
	UserID    uint32
}

var (
	err      error
	bodyRaw  []byte
	req      *http.Request
	countRec int
)

func (s *mySuite) CheckCountRec(expCount int, myQueryText string) {
	mySQLRows, err := s.DBConnect.QueryContext(s.ctx, myQueryText)
	s.Require().NoError(err)
	defer mySQLRows.Close()
	mySQLRows.Next()
	err = mySQLRows.Scan(&countRec)
	fmt.Println("countRec=", countRec)
	s.Require().NoError(err)
	s.Require().Equal(expCount, countRec)
}

func (s *mySuite) CheckCountRecBind(expCount int, myQueryText string, arg interface{}) {
	fmt.Println("args=", arg, myQueryText)
	mySQLRows, err := s.DBConnect.QueryContext(s.ctx, myQueryText, arg)
	s.Require().NoError(err)
	defer mySQLRows.Close()
	mySQLRows.Next()
	err = mySQLRows.Scan(&countRec)
	fmt.Println("countRec=", countRec)
	s.Require().NoError(err)
	s.Require().Equal(expCount, countRec)
}

func (s *mySuite) GetMaxValBind(myQueryText string, arg interface{}) string {
	var maxVal string
	mySQLRows, err := s.DBConnect.QueryContext(s.ctx, myQueryText, arg)
	s.Require().NoError(err)
	defer mySQLRows.Close()
	mySQLRows.Next()
	err = mySQLRows.Scan(&maxVal)
	s.Require().NoError(err)
	return maxVal
}

func (s *mySuite) SetupSuite() {
	s.client = http.Client{
		Timeout: time.Second * 5,
	}

	s.hostName = "http://mainSevice:5000/" // через докер
	//	s.hostName = "http://localhost:5000/" // локально
	s.ctx = context.Background()

	myStr := "postgres://calend:calend@postgres_db:5432/calend?sslmode=disable" // через докер
	//	myStr := "postgres://calend:calend@localhost:5432/calend?sslmode=disable" // локально

	s.DBConnect, err = sql.Open("postgres", myStr)
	if err == nil {
		err = s.DBConnect.PingContext(s.ctx)
	}
	s.Require().NoError(err)

	_, err = s.DBConnect.ExecContext(s.ctx, "delete from events")
	s.Require().NoError(err)
	s.CheckCountRec(0, "select count(*) RC from events")

	_, err = s.DBConnect.ExecContext(s.ctx, "delete from send_events_stat")
	s.Require().NoError(err)
	s.CheckCountRec(0, "select count(*) RC from send_events_stat")

	_, err = s.DBConnect.ExecContext(s.ctx, "delete from shed_send_id")
	s.Require().NoError(err)
	s.CheckCountRec(0, "select count(*) RC from send_events_stat")
}

func (s *mySuite) TearDownSuite() {
	mySQL := `
				delete from events;
				delete from shed_send_id;
				delete from send_events_stat;
			`
	_, err = s.DBConnect.ExecContext(s.ctx, mySQL)
	s.Require().NoError(err)
	s.DBConnect.Close()
}

func (s *mySuite) SendRequest(myMethodName string, myAnyStruct interface{}) []byte {
	bodyRaw, err = json.Marshal(myAnyStruct)
	s.Require().NoError(err)
	req, err = http.NewRequestWithContext(s.ctx, http.MethodPost, s.hostName+myMethodName, bytes.NewBuffer(bodyRaw))
	s.Require().NoError(err)

	newResp, err2 := s.client.Do(req)
	s.Require().NoError(err2)

	defer newResp.Body.Close()
	bodyRaw, err = ioutil.ReadAll(newResp.Body)
	s.Require().NoError(err)

	return bodyRaw
}

func (s *mySuite) CreateEvent(myForCreate ForCreate) {
	bodyRaw = s.SendRequest("CreateEvent", myForCreate)
	s.Require().Empty(bodyRaw)
}

func (s *mySuite) UpdateEvent(myForUpdate ForUpdate) {
	bodyRaw = s.SendRequest("UpdateEvent", myForUpdate)
	s.Require().Empty(bodyRaw)
}

func (s *mySuite) DeleteEvent(myForDelete ForDelete) {
	bodyRaw = s.SendRequest("DeleteEvent", myForDelete)
	s.Require().Empty(bodyRaw)
}

func (s *mySuite) GetEventByDate(myStartDate StartDate) []OneEvent {
	var myEventList []OneEvent
	bodyRaw = s.SendRequest("GetEventByDate", myStartDate)
	s.Require().NotEmpty(bodyRaw)
	err = json.Unmarshal(bodyRaw, &myEventList)
	s.Require().NoError(err)
	return myEventList
}

func (s *mySuite) GetEventWeek(myStartDate StartDate, expectNil bool) []OneEvent {
	var myEventList []OneEvent
	bodyRaw = s.SendRequest("GetEventWeek", myStartDate)
	if expectNil {
		s.Require().Empty(bodyRaw)
		return nil
	}
	s.Require().NotEmpty(bodyRaw)
	err = json.Unmarshal(bodyRaw, &myEventList)
	s.Require().NoError(err)
	return myEventList
}

func (s *mySuite) GetEventMonth(myStartDate StartDate, expectNil bool) []OneEvent {
	var myEventList []OneEvent
	bodyRaw = s.SendRequest("GetEventMonth", myStartDate)
	if expectNil {
		s.Require().Empty(bodyRaw)
		return nil
	}
	s.Require().NotEmpty(bodyRaw)
	err = json.Unmarshal(bodyRaw, &myEventList)
	s.Require().NoError(err)
	return myEventList
}

func (s *mySuite) Test1AddDelUpd() {
	fmt.Println("Start Test1")
	constDate := "21.01.2023"

	mySQL := "select count(*) RC from events where StartDate=to_date($1,'DD.MM.YYYY')"
	mySQL2 := "select max(Title) MV from events where ID=$1"
	s.CheckCountRecBind(0, mySQL, constDate)
	s.CreateEvent(ForCreate{Title: "newTitle", StartDate: constDate, Details: "newDetails", UserID: 123})
	s.CheckCountRecBind(1, mySQL, constDate)
	newID := s.GetMaxValBind("select max(ID) MV from events where StartDate=to_date($1,'DD.MM.YYYY')", constDate)

	myTitle := s.GetMaxValBind(mySQL2, newID)
	s.Require().Equal("newTitle", myTitle)
	s.UpdateEvent(ForUpdate{EventID: newID, Title: "newTitle99", StartDate: constDate, Details: "newDetails", UserID: 123})
	myTitle = s.GetMaxValBind(mySQL2, newID)
	s.Require().Equal("newTitle99", myTitle)

	s.CreateEvent(ForCreate{Title: "newTitle2", StartDate: constDate, Details: "newDetails2", UserID: 124})
	s.CheckCountRecBind(2, mySQL, constDate)

	s.DeleteEvent(ForDelete{EventID: newID})
	s.CheckCountRecBind(1, mySQL, constDate)
	fmt.Println("finish Test1AddBanner")
}

func (s *mySuite) Test2GetEvent() {
	fmt.Println("Start Test2")
	var myEventList []OneEvent
	constDate := "05.12.2022"
	constDate2 := "01.12.2022"
	constDate3 := "06.12.2022"
	constDate4 := "15.12.2022"

	myEventList = s.GetEventWeek(StartDate{StartDateStr: constDate2}, true) // не начало недели
	s.Require().Empty(myEventList)
	myEventList = s.GetEventMonth(StartDate{StartDateStr: constDate}, true) // не начало месяца
	mySQL := "select count(*) RC from events where StartDate=to_date($1,'DD.MM.YYYY')"
	s.Require().Empty(myEventList)

	s.CheckCountRecBind(0, mySQL, constDate)
	myEventList = s.GetEventByDate(StartDate{StartDateStr: constDate})
	s.Require().Equal(0, len(myEventList))
	myEventList = s.GetEventWeek(StartDate{StartDateStr: constDate}, false)
	s.Require().Equal(0, len(myEventList))
	myEventList = s.GetEventMonth(StartDate{StartDateStr: constDate2}, false)
	s.Require().Equal(0, len(myEventList))

	myForCreate := ForCreate{Title: "newTitle", StartDate: constDate, Details: "newDetails", UserID: 123}
	s.CreateEvent(myForCreate)
	s.CheckCountRecBind(1, mySQL, constDate)
	myEventList = s.GetEventByDate(StartDate{StartDateStr: constDate})
	s.Require().Equal(1, len(myEventList))
	s.Require().Equal(myForCreate.Title, myEventList[0].Title)
	s.Require().Equal(myForCreate.Details, myEventList[0].Details)
	s.Require().Equal(myForCreate.UserID, myEventList[0].UserID)

	s.CreateEvent(ForCreate{Title: "newTitle+", StartDate: constDate3, Details: "newDetails2", UserID: 1})
	s.CheckCountRecBind(1, mySQL, constDate3)
	s.CreateEvent(ForCreate{Title: "newTitle-", StartDate: constDate4, Details: "newDetails3", UserID: 12355})
	s.CheckCountRecBind(1, mySQL, constDate4)

	myEventList = s.GetEventByDate(StartDate{StartDateStr: constDate})
	s.Require().Equal(1, len(myEventList))
	myEventList = s.GetEventWeek(StartDate{StartDateStr: constDate}, false)
	s.Require().Equal(2, len(myEventList))
	myEventList = s.GetEventMonth(StartDate{StartDateStr: constDate2}, false)
	s.Require().Equal(3, len(myEventList))
	fmt.Println("finish Test2")
}

func (s *mySuite) Test3SchedulerSender() {
	fmt.Println("Start Test1")
	constDate := time.Now().Format("02.01.2006")
	mySQL := "select count(*) RC from events where StartDate=to_date($1,'DD.MM.YYYY')"
	mySQL2 := "select count(*) RC from shed_send_id"
	mySQL3 := "select count(*) RC from send_events_stat"
	s.CheckCountRecBind(0, mySQL, constDate)
	s.CheckCountRec(0, mySQL2)
	s.CheckCountRec(0, mySQL3)

	s.CreateEvent(ForCreate{Title: "newTitle", StartDate: constDate, Details: "newDetails", UserID: 1})
	s.CheckCountRecBind(1, mySQL, constDate)
	s.CheckCountRec(0, mySQL2)
	s.CheckCountRec(0, mySQL3)

	time.Sleep(time.Second * 35)
	s.CheckCountRec(1, mySQL2)
	s.CheckCountRec(1, mySQL3)
}

func TestService(t *testing.T) {
	suite.Run(t, new(mySuite))
}

/*
func (s *mySuite) GetBannerForSlot(mySlotSoc ForGetBanner) GetBannerStruct { // получения баннера для показа в слоте
	bodyRaw = s.SendRequest("GetBannerForSlot", mySlotSoc)
	s.Require().NotEmpty(bodyRaw)
	err = json.Unmarshal(bodyRaw, &myGetBannerStruct)
	s.Require().NoError(err)
	return myGetBannerStruct
}

func (s *mySuite) BannerClick(myBannerClick ForBannerClick) { // клик по баннеру
	bodyRaw = s.SendRequest("BannerClick", myBannerClick)
	s.Require().Empty(bodyRaw)
}

func (s *mySuite) Test1AddBanner() {
	for i := 1; i <= 10; i++ {
		s.AddSlotBanner(SlotBanner{1, i})
	}
	// к слоту привязано 10 баннеров
	s.CheckCountRec("select count(*) RC from slot_banner where slot_id=1", 10)
	// привязан баннер с id=1
	s.CheckCountRec("select count(*) RC from slot_banner where slot_id=1 and banner_id=1", 1)
	s.AddSlotBanner(SlotBanner{1, 1})
	// после повторной попытке ничего не изменилось
	s.CheckCountRec("select count(*) RC from slot_banner where slot_id=1 and banner_id=1", 1)
	fmt.Println("finish Test1AddBanner")
}

func (s *mySuite) Test2DelBanner() {
	fmt.Println("start TestDelSlotBanner")
	//  добавим баннер к слоту
	s.AddSlotBanner(SlotBanner{2, 2})
	// убедимся что он добавился
	s.CheckCountRec("select count(*) RC from slot_banner where slot_id=2 and banner_id=2", 1)
	// отвяжем баннер от слота
	s.DelSlotBanner(SlotBanner{2, 2})
	// убедимся что отвязался
	s.CheckCountRec("select count(*) RC from slot_banner where slot_id=2 and banner_id=2", 0)
	fmt.Println("finish TestDelSlotBanner")
}

func (s *mySuite) Test3GetBannerForSlot() {
	// убедимся, что к слоту 2 не првязан ни один баннер
	s.CheckCountRec("select count(*) RC from slot_banner where slot_id=2", 0)
	// поскольку к слоту 2 не првязан ни один баннер должен вернуть 0
	myGetBannerStruct = s.GetBannerForSlot(ForGetBanner{2, 1})
	s.Require().Equal(GetBannerStruct{}, myGetBannerStruct)

	// добавим во второй слот баннер с ID=3
	s.AddSlotBanner(SlotBanner{2, 3})
	// теперь должен вернуть ID=3 (так как это единственный баннер
	myGetBannerStruct = s.GetBannerForSlot(ForGetBanner{2, 1})
	s.Require().Equal(3, myGetBannerStruct.ID)
	// убедимся что этот показ отразился в статистике (1 раз)
	s.CheckCountRec("select count(*) RC from banner_stat where stat_type='S' and slot_id=2 and banner_id=3", 1)
}

func (s *mySuite) Test4BannerClick() {
	// убедимся, что к в слоте 1 для баннера 2 для соц группы 3 ещё не было кликов
	mySQL := `
	  select count(*) RC
	  from banner_stat
	  where stat_type='C'
		and slot_id=1
		and banner_id=2
		and soc_group_id=3`
	s.CheckCountRec(mySQL, 0)
	//  кликнем в слоте 1 на баннер 2 для соц группы 3
	s.BannerClick(ForBannerClick{1, 2, 3})
	// убедимся, что теперь сохранился 1 клик
	s.CheckCountRec(mySQL, 1)
}

func (s *mySuite) Test5SendMessages() {
	mySQL := `
	  select count(*) RC
      from banner_stat
	  where id>(
		select max(banner_stat_id)
		from send_stat_max_id
	)`
	// убедимся, что есть неотправленные сообщения (2)
	s.CheckCountRec(mySQL, 2)

	// уснём, чтобы дождаться отправки статистики
	time.Sleep(time.Second * 30)

	// убедимся, что теперь не осталоьс не отправленных сообщений
	s.CheckCountRec(mySQL, 0)
}
*/

/*
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
*/
