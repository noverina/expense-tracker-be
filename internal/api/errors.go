package api

import "errors"

var ErrAtoiFail = errors.New("unable to convert from string to integer")
var ErrInvalidType = errors.New("invalid type")
var ErrInvalidCategory = errors.New("invalid category")
var ErrAmountInvalidChar = errors.New("amount must only contain numbers")
var ErrMaxEvent = errors.New("event limit")
var ErrMaxMonthRange = errors.New("month range limit")
var ErrInvalidFilterKey = errors.New("invalid field to query")
var ErrInvalidToken = errors.New("jwt is invalid")
