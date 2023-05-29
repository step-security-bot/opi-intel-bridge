// SPDX-License-Identifier: Apache-2.0
// Copyright (C) 2023 Intel Corporation

// Package utils implements utility functions
package utils

import "google.golang.org/protobuf/proto"

// ProtoClone creates a deep copy and type asserts to a provided instance
func ProtoClone[T proto.Message](protoStruct T) T {
	return proto.Clone(protoStruct).(T)
}
