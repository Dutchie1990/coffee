package helpers_test

import (
	"bytes"
	"coffee/coffee-server/helpers"
	"coffee/coffee-server/services"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helpers", Label("unit"), func() {
	var (
		w       *httptest.ResponseRecorder
		request *http.Request
	)

	BeforeEach(func() {
		w = httptest.NewRecorder()
	})

	Describe("ReadJson", func() {
		type TestPayload struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}

		Context("with valid JSON", func() {
			It("should parse the JSON body correctly", func() {
				payload := `{"name": "John Doe", "email": "johndoe@example.com"}`
				request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(payload)))

				var data TestPayload
				err := helpers.ReadJson(w, request, &data)

				Expect(err).To(BeNil())
				Expect(data.Name).To(Equal("John Doe"))
				Expect(data.Email).To(Equal("johndoe@example.com"))
			})
		})

		Context("with invalid JSON", func() {
			It("should return an error for malformed JSON", func() {
				payload := `{"name": "John Doe", "email":}`
				request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(payload)))

				var data TestPayload
				err := helpers.ReadJson(w, request, &data)

				Expect(err).ToNot(BeNil())
			})
		})

		Context("when the body has multiple JSON objects", func() {
			It("should return an error when multiple JSON objects are present", func() {
				payload := `{"name": "John Doe"} {"email": "johndoe@example.com"}`
				request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(payload)))

				var data TestPayload
				err := helpers.ReadJson(w, request, &data)

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("body must have only a single JSON object"))
			})
		})
	})

	Describe("WriteJson", func() {
		Context("with valid data", func() {
			It("should write a JSON response with correct headers and status", func() {
				data := map[string]string{
					"message": "Success",
				}

				err := helpers.WriteJson(w, http.StatusOK, data)

				Expect(err).To(BeNil())
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Header().Get("Content-Type")).To(Equal("application/json"))

				var response map[string]string
				err = json.NewDecoder(w.Body).Decode(&response)

				Expect(err).To(BeNil())
				Expect(response["message"]).To(Equal("Success"))
			})
		})
	})

	Describe("ErrorJson", func() {
		Context("when an error is passed", func() {
			It("should return a JSON error response with the correct status", func() {
				err := errors.New("an error occurred")
				helpers.ErrorJson(w, err)

				Expect(w.Code).To(Equal(http.StatusBadRequest))

				var response services.JsonResponse
				err = json.NewDecoder(w.Body).Decode(&response)

				Expect(err).To(BeNil())
				Expect(response.Error).To(BeTrue())
				Expect(response.Message).To(Equal("an error occurred"))
			})

			It("should allow passing a custom status code", func() {
				err := errors.New("not found")
				helpers.ErrorJson(w, err, http.StatusNotFound)

				Expect(w.Code).To(Equal(http.StatusNotFound))

				var response services.JsonResponse
				err = json.NewDecoder(w.Body).Decode(&response)

				Expect(err).To(BeNil())
				Expect(response.Error).To(BeTrue())
				Expect(response.Message).To(Equal("not found"))
			})
		})
	})
})
