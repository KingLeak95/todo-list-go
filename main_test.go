package main
import (
  "github.com/gin-gonic/gin"
  "github.com/DATA-DOG/go-sqlmock"
  "testing"
)

func SetUpRouter() *gin.Engine{
    router := gin.Default()
    return router
}

func TestCreateUser(t *testing.T) {
    db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
    r := SetUpRouter()
    r.POST("/company", NewCompanyHandler)
    companyId := xid.New().String()
    company := Company{
        ID: companyId,
        Name: "Demo Company",
        CEO: "Demo CEO",
        Revenue: "35 million",
    }
    jsonValue, _ := json.Marshal(company)
    req, _ := http.NewRequest("POST", "/company", bytes.NewBuffer(jsonValue))

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusCreated, w.Code)
}
