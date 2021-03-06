package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nelsonmhjr/bank_service/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Accounts", func() {
	AfterEach(func() {
		Db.Exec("Truncate bank_accounts")
	})

	var requestBody = func(resp *http.Response) string {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		Expect(err).To(BeNil())
		return string(bodyBytes)
	}

	Describe("FindAccount GET /accounts/:accountId", func() {
		var request = func(id uint) *http.Response {
			resp, err := http.Get(fmt.Sprintf("%s/accounts/%d", Ts.URL, id))
			Expect(err).To(BeNil())
			return resp
		}

		Context("without a matching BankAccount with the matchin ID", func() {
			It("should render a Bad Request", func() {
				resp := request(1)
				Expect(resp.StatusCode).To(Equal(400))
			})
		})

		Context("with a matching BankAccount ID", func() {
			It("should render JSON with id and document_number", func() {
				bankAcc := models.BankAccount{DocumentNumber: "123456"}
				Db.Create(&bankAcc)
				resp := request(bankAcc.ID)
				bodyString := requestBody(resp)
				Expect(bodyString).To(
					MatchJSON(
						fmt.Sprintf(`{"document_number":"%s","account_id":%d}`,
							bankAcc.DocumentNumber, bankAcc.ID)))
			})
		})
	})

	Describe("CreateAccounts POST /accounts", func() {
		var request = func(accToCreate models.AccountToCreate) *http.Response {
			body, err := json.Marshal(accToCreate)
			Expect(err).To(BeNil())
			resp, err := http.Post(fmt.Sprintf("%s/accounts/", Ts.URL),
				"application/json", bytes.NewBuffer(body))
			Expect(err).To(BeNil())
			return resp
		}

		Context("without a json body", func() {
			It("should render a UnprocessableEntity Status Code", func() {
				resp := request(models.AccountToCreate{})
				Expect(resp.StatusCode).To(Equal(422))
			})
		})

		Context("with a valid json body", func() {
			It("should create the account", func() {
				var cntBef, cntAft int
				Db.Model(&models.BankAccount{}).Count(&cntBef)
				resp := request(models.AccountToCreate{
					DocumentNumber: "12345678"})

				Expect(resp.StatusCode).To(Equal(201))
				Db.Model(&models.BankAccount{}).Count(&cntAft)
				Expect(cntAft).To(Equal(cntBef + 1))
			})
		})
	})
})
