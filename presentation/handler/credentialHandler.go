package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/miliya612/webauthn-demo/presentation/httputil"
	"github.com/miliya612/webauthn-demo/presentation/usecase"
	"github.com/miliya612/webauthn-demo/presentation/usecase/input"
	"io"
	"io/ioutil"
	"net/http"
)

var KB int64 = 1024

type CredentialHandler interface {
	RegistrationInit(w http.ResponseWriter, r *http.Request)
	//Registration(w http.ResponseWriter, r *http.Request)
	//AuthenticationInit(w http.ResponseWriter, r *http.Request)
	//Authentication(w http.ResponseWriter, r *http.Request)
}

type credentialHandler struct {
	registrationInit usecase.RegistrationInitUseCase
	registration usecase.RegistrationUseCase
}

func NewCredentialHandler(
		registrationInit usecase.RegistrationInitUseCase,
		//registration usecase.RegistrationUseCase,
	) CredentialHandler {
	return &credentialHandler{
		registrationInit: registrationInit,
		//registration: registration,
	}
}

func (h *credentialHandler) RegistrationInit(w http.ResponseWriter, r *http.Request) {
	in, err := parseRegistrationInitRequest(r)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "failed to parse request", err)
		return
	}

	if err = validateRegistrationInitRequest(in); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request", err)
		return
	}

	resp, err := h.registrationInit.RegistrationInit(*in)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}
	httputil.Ok(w, resp)
}


func (h *credentialHandler) Registration(w http.ResponseWriter, r *http.Request) {
	var in input.Registration
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*KB)) // 1MB
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "request body is too large", err)
	}
	defer r.Body.Close()

	if err = json.Unmarshal(body, &in); err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed marshalling json", err)
	}

	resp, err := h.registration.Registration(in)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "something went wrong", err)
	}
	httputil.Ok(w, resp)
}

func parseRegistrationInitRequest(r *http.Request) (*input.RegistrationInit, error) {
	var in input.RegistrationInit
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*KB)) // 1MB
	if err != nil {
		return nil, errors.New("request body is too large")
	}
	defer r.Body.Close()

	if len(body) == 0 {
		return nil, errors.New("empty body is not acceptable")
	}

	if err = json.Unmarshal(body, &in); err != nil {
		return nil, errors.New(fmt.Sprint("failed marshalling json", err))
	}


	return &in, nil
}

func validateRegistrationInitRequest(in *input.RegistrationInit) error {
	var invalidParams []string
	if in.ID == "" {
		invalidParams = append(invalidParams, "id")
	}
	if in.DisplayName == "" {
		invalidParams = append(invalidParams, "displayName")
	}
	if len(invalidParams) != 0 {
		return errors.New(fmt.Sprint("required parameters are missing: ", invalidParams))
	}
	return nil
}