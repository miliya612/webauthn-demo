package handler

import (
	"encoding/json"
	"github.com/miliya612/webauthn-demo/presentation/httputil"
	"github.com/miliya612/webauthn-demo/presentation/usecase"
	"github.com/miliya612/webauthn-demo/presentation/usecase/input"
	"io"
	"io/ioutil"
	"net/http"
)

var KB int64 = 1024

type CredentialHandler interface {
	RegistrationInit(r *http.Request) httputil.Responder
	//Registration(r *http.Request) httputil.Responder
	//AuthenticationInit(r *http.Request) httputil.Responder
	//Authentication(r *http.Request) httputil.Responder
}

type credentialHandler struct {
	registrationInit usecase.RegistrationInitUseCase
	registration usecase.RegistrationUseCase
}

func NewCredentialHandler(
		registrationInit usecase.RegistrationInitUseCase,
		registration usecase.RegistrationUseCase,
	) CredentialHandler {
	return &credentialHandler{
		registrationInit: registrationInit,
		registration: registration,
	}
}

func (h *credentialHandler) RegistrationInit(r *http.Request) httputil.Responder {
	var in input.RegistrationInit
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*KB)) // 1MB
	if err != nil {
		return httputil.Error(http.StatusInternalServerError, "request body is too large", err)
	}
	defer r.Body.Close()

	if err = json.Unmarshal(body, &in); err != nil {
		return httputil.Error(http.StatusInternalServerError, "failed marshalling json", err)
	}

	resp, err := h.registrationInit.RegistrationInit(in)
	if err != nil {
		return httputil.Error(http.StatusInternalServerError, "something went wrong", err)
	}
	return httputil.Ok(resp)
}


func (h *credentialHandler) Registration(r *http.Request) httputil.Responder {
	var in input.Registration
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*KB)) // 1MB
	if err != nil {
		return httputil.Error(http.StatusInternalServerError, "request body is too large", err)
	}
	defer r.Body.Close()

	if err = json.Unmarshal(body, &in); err != nil {
		return httputil.Error(http.StatusInternalServerError, "failed marshalling json", err)
	}

	resp, err := h.registration.Registration(in)
	if err != nil {
		return httputil.Error(http.StatusInternalServerError, "something went wrong", err)
	}
	return httputil.Ok(resp)
}

//
//func parseCredentialId(r *http.Request) (int, error) {
//	id, err := strconv.Atoi(mux.Vars(r)["credentialId"])
//	if err != nil {
//		return -1, errors.New("credentialId should be number.")
//	}
//	return id, nil
//}
