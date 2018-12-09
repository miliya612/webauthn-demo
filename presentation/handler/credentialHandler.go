package handler

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/miliya612/webauthn-demo/presentation/httputil"
	"github.com/miliya612/webauthn-demo/presentation/usecase"
	"github.com/miliya612/webauthn-demo/presentation/usecase/input"
	"net/http"
)

type CredentialHandler interface {
	RegistrationInit(w http.ResponseWriter, r *http.Request)
	Registration(w http.ResponseWriter, r *http.Request)
	//AuthenticationInit(w http.ResponseWriter, r *http.Request)
	//Authentication(w http.ResponseWriter, r *http.Request)
}

type credentialHandler struct {
	registrationInit usecase.RegistrationInitUseCase
	registration     usecase.RegistrationUseCase
}

func NewCredentialHandler(
	registrationInit usecase.RegistrationInitUseCase,
	registration usecase.RegistrationUseCase,
) CredentialHandler {
	return &credentialHandler{
		registrationInit: registrationInit,
		registration:     registration,
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

	ctx := r.Context()
	uuid, err := createUUID()
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}
	ctx = context.WithValue(ctx, "sid", uuid)

	resp, err := h.registrationInit.RegistrationInit(ctx, *in)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	fmt.Println("chal: ", string(resp.PublicKey.Challenge))

	http.SetCookie(w, &http.Cookie{Name: "sid", Value: uuid})

	httputil.Created(w, resp)
}

func (h *credentialHandler) Registration(w http.ResponseWriter, r *http.Request) {
	in, err := parseRegistrationRequest(r)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "failed to parse request", err)
		return
	}

	if err = validateRegistrationRequest(in); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request", err)
		return
	}

	resp, err := h.registration.Registration(r.Context(), *in)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}
	httputil.Created(w, resp)
}

func parseRegistrationInitRequest(r *http.Request) (*input.RegistrationInit, error) {
	var in input.RegistrationInit
	body, err := httputil.ParseBody(r)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &in); err != nil {
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

func parseRegistrationRequest(r *http.Request) (*input.Registration, error) {
	var in input.Registration

	body, err := httputil.ParseBody(r)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &in); err != nil {
		return nil, errors.New(fmt.Sprint("failed marshalling json", err))
	}

	c, err := r.Cookie("sid")
	if err != nil {
		return nil, errors.New(fmt.Sprint("failed parsing cookie", err))
	}

	fmt.Println("sessionID: ", c.Value)
	ctx := context.WithValue(r.Context(), "sid", c.Value)
	r = r.WithContext(ctx)

	return &in, nil
}

func validateRegistrationRequest(in *input.Registration) error {
	var invalidParams []string
	//if in.ID == "" {
	//	invalidParams = append(invalidParams, "id")
	//}
	//if in.DisplayName == "" {
	//	invalidParams = append(invalidParams, "displayName")
	//}
	if len(invalidParams) != 0 {
		return errors.New(fmt.Sprint("required parameters are missing: ", invalidParams))
	}
	return nil
}

// create a random UUID with from RFC 4122
// adapted from http://github.com/nu7hatch/gouuid
func createUUID() (string, error){
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to generate UUID: %s", err))
	}

	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:]), nil
}
