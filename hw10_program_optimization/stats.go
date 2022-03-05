package hw10programoptimization

// Оптимизация по времени далась сравнительно легко (4 оптимизации), а вот по памяти...
// Дошёл до 45 Мб, и дальше не удавалось сдвинуться (даже попытка использование sync.Pool ничего не дало)...
// Единственным вариантом для меня, стало отойти от использования Unmarshal (зря только подключил  Easy Json),
// создать минимальную структуру и использовать регулярку (от которой ушёл в другом месте) для получения данных...
// Вроде условия задания не были нарушены, но я на всякий случай не стал удалять предыдущее решение...
import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type User2 struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	//	u, err := getUsers(r)
	if domain == "" {
		return nil, fmt.Errorf("domain is empty")
	}

	u, err := getUsers2(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type (
	//	users  [100_000]User
	users2 [100_000]User2
)

func getUsers2(r io.Reader) (result users2, err error) {
	i := 0
	myPattern := `"Email":"(.*?)"`
	myReg, err := regexp.Compile(myPattern)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		myEmail := myReg.FindStringSubmatch(line)
		result[i].Email = myEmail[1]
		i++
	}
	return
}

func countDomains(u /*users*/ users2, domain string) (DomainStat, error) {
	result := make(DomainStat)
	myDomain := "." + domain
	for _, user := range u {
		matched := strings.Contains(user.Email, myDomain)
		if matched {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result, nil
}

/*func getUsers(r io.Reader) (result users, err error) {
	i := 0
	//	p := sync.Pool{
	//		New: func() interface{} {
	//			return &User{}
	//		},
	//	}
	var user User
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		//		user := p.Get().(*User)
		// if err = json.Unmarshal([]byte(line), &user); err != nil {
		if err = user.UnmarshalJSON([]byte(line)); err != nil {
			return
		}
		// result[i] = *user
		// *user = User{}
		// p.Put(user)
		result[i] = user
		i++
	}
	return
}*/
