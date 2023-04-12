// SPDX-License-Identifier: Apache-2.0
// Copyright (C) 2023 Intel Corporation

// Package middleend implements the MiddleEnd APIs (service) of the storage Server
package middleend

import (
	"context"
	"encoding/hex"
	"log"
	"runtime"
	"runtime/debug"
	"strings"

	pc "github.com/opiproject/opi-api/common/v1/gen/go"
	pb "github.com/opiproject/opi-api/storage/v1alpha1/gen/go"
	intelmodels "github.com/opiproject/opi-intel-bridge/pkg/models"
	"github.com/opiproject/opi-spdk-bridge/pkg/middleend"
	spdkmodels "github.com/opiproject/opi-spdk-bridge/pkg/models"
	"github.com/opiproject/opi-spdk-bridge/pkg/server"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	errMissingArgument    = status.Error(codes.InvalidArgument, "missing argument")
	errNotSupportedCipher = status.Error(codes.Unimplemented, "not supported cipher")
	errWrongKeySize       = status.Error(codes.InvalidArgument, "invalid key size")
)

// Server contains middleend related OPI services
type Server struct {
	pb.MiddleendServiceServer
	rpc server.JSONRPC
}

// NewServer creates initialized instance of middleend server
func NewServer(jsonRPC server.JSONRPC) *Server {
	opiSpdkServer := middleend.NewServer(jsonRPC)
	return &Server{
		opiSpdkServer,
		jsonRPC,
	}
}

// CreateEncryptedVolume creates an encrypted volume
func (s *Server) CreateEncryptedVolume(_ context.Context, in *pb.CreateEncryptedVolumeRequest) (*pb.EncryptedVolume, error) {
	log.Printf("CreateEncryptedVolume was called")
	if in == nil {
		log.Println("request cannot be empty")
		return nil, errMissingArgument
	}
	defer cleanupEncryptedVolumeStruct(in.EncryptedVolume)
	if err := verifyEncryptedVolume(in.EncryptedVolume); err != nil {
		return nil, err
	}

	bdevUUID, err := s.getBdevUUIDByName(in.EncryptedVolume.VolumeId.Value)
	if err != nil {
		log.Println("Failed to find UUID for bdev", in.EncryptedVolume.VolumeId.Value)
		return nil, err
	}

	half := len(in.EncryptedVolume.Key) / 2
	tweakMode := "A"
	params := &intelmodels.NpiBdevSetKeysParams{
		UUID:   bdevUUID,
		Key:    hex.EncodeToString(in.EncryptedVolume.Key[:half]),
		Key2:   hex.EncodeToString(in.EncryptedVolume.Key[half:]),
		Cipher: "AES_XTS",
		Tweak:  tweakMode,
	}
	defer func() {
		params = nil
	}()

	var result intelmodels.NpiBdevSetKeysResult
	err = s.rpc.Call("npi_bdev_set_keys", params, &result)
	if err != nil {
		log.Println("error:", err)
		return nil, server.ErrFailedSpdkCall
	}
	if !result {
		log.Println("Failed result on SPDK call:", result)
		return nil, server.ErrUnexpectedSpdkCallResult
	}

	return &pb.EncryptedVolume{
		EncryptedVolumeId: &pc.ObjectKey{Value: in.EncryptedVolume.EncryptedVolumeId.Value},
		VolumeId:          &pc.ObjectKey{Value: in.EncryptedVolume.VolumeId.Value},
	}, nil
}

// DeleteEncryptedVolume deletes an encrypted volume
func (s *Server) DeleteEncryptedVolume(_ context.Context, in *pb.DeleteEncryptedVolumeRequest) (*emptypb.Empty, error) {
	log.Printf("DeleteEncryptedVolume: Received from client: %v", in)
	if in == nil {
		log.Println("request cannot be empty")
		return nil, errMissingArgument
	}

	bdevUUID, err := s.getBdevUUIDByName(in.Name)
	if err != nil {
		log.Println("Failed to find UUID for bdev", in.Name)
		return nil, err
	}
	params := intelmodels.NpiBdevClearKeysParams{
		UUID: bdevUUID,
	}
	var result intelmodels.NpiBdevClearKeysResult
	err = s.rpc.Call("npi_bdev_clear_keys", params, &result)
	if err != nil {
		cryptoObjMissingErrMsg := "Could not find a crypto object for a given bdev"
		if in.AllowMissing &&
			strings.Contains(err.Error(), cryptoObjMissingErrMsg) {
			return &emptypb.Empty{}, nil
		}
		log.Println(err)
		return nil, server.ErrFailedSpdkCall
	}
	if !result {
		log.Println("Failed result on SPDK call:", result)
		return nil, server.ErrUnexpectedSpdkCallResult
	}
	return &emptypb.Empty{}, nil
}

// UpdateEncryptedVolume updates an encrypted volume
func (s *Server) UpdateEncryptedVolume(ctx context.Context, in *pb.UpdateEncryptedVolumeRequest) (*pb.EncryptedVolume, error) {
	log.Printf("UpdateEncryptedVolume was called")
	if in == nil {
		log.Println("request cannot be empty")
		return nil, errMissingArgument
	}
	defer cleanupEncryptedVolumeStruct(in.EncryptedVolume)
	if err := verifyEncryptedVolume(in.EncryptedVolume); err != nil {
		return nil, err
	}

	if _, err := s.DeleteEncryptedVolume(ctx, &pb.DeleteEncryptedVolumeRequest{
		Name: in.EncryptedVolume.VolumeId.Value,
	}); err != nil {
		return nil, err
	}

	return s.CreateEncryptedVolume(ctx, &pb.CreateEncryptedVolumeRequest{
		EncryptedVolume: in.EncryptedVolume,
	})
}

func (s *Server) getBdevUUIDByName(name string) (string, error) {
	params := spdkmodels.BdevGetBdevsParams{Name: name}
	var result []spdkmodels.BdevGetBdevsResult
	err := s.rpc.Call("bdev_get_bdevs", params, &result)
	if err != nil {
		log.Println("error:", err)
		return "", server.ErrFailedSpdkCall
	}
	if len(result) != 1 {
		log.Println("Found bdevs:", result, "under the name", params.Name)
		return "", server.ErrUnexpectedSpdkCallResult
	}
	return result[0].UUID, nil
}

func verifyEncryptedVolume(volume *pb.EncryptedVolume) error {
	switch {
	case volume == nil:
		log.Println("encrypted_volume should be specified")
		return errMissingArgument
	case volume.EncryptedVolumeId == nil || volume.EncryptedVolumeId.Value == "":
		log.Println("encrypted_volume_id should be specified")
		return errMissingArgument
	case volume.VolumeId == nil || volume.VolumeId.Value == "":
		log.Println("volume_id should be specified")
		return errMissingArgument
	case len(volume.Key) == 0:
		log.Println("key cannot be empty")
		return errMissingArgument
	}

	keyLengthInBits := len(volume.Key) * 8
	expectedKeyLengthInBits := 0
	switch {
	case volume.Cipher == pb.EncryptionType_ENCRYPTION_TYPE_AES_XTS_256:
		expectedKeyLengthInBits = 512
	case volume.Cipher == pb.EncryptionType_ENCRYPTION_TYPE_AES_XTS_128:
		expectedKeyLengthInBits = 256
	default:
		log.Println("only AES_XTS_128 and AES_XTS_256 are supported")
		return errNotSupportedCipher
	}

	if keyLengthInBits != expectedKeyLengthInBits {
		log.Printf("expected key size %vb, provided size %vb",
			expectedKeyLengthInBits, keyLengthInBits)
		return errWrongKeySize
	}

	return nil
}

func cleanupEncryptedVolumeStruct(volume *pb.EncryptedVolume) {
	if volume != nil {
		for i := range volume.Key {
			volume.Key[i] = 0
		}
		volume.Cipher = pb.EncryptionType_ENCRYPTION_TYPE_UNSPECIFIED
	}
	// Run GC to free all variables which contained encryption keys
	runtime.GC()
	// Return allocated memory to OS, otherwise the memory which contained
	// keys can be kept for some time.
	debug.FreeOSMemory()
}
