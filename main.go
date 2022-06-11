package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Clarifai/clarifai-go-grpc/proto/clarifai/api"
	"github.com/Clarifai/clarifai-go-grpc/proto/clarifai/api/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var API_KEY = os.Getenv("API_KEY")

func main() {
	// NOTE: ここでclarifaiのAPIに対してgRPC通信を用いてコネクションを繋いでいる。
	conn, err := grpc.Dial(
		"api.clarifai.com:443",
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
	)
	if err != nil {
		panic(err)
	}
	client := api.NewV2Client(conn)

	// NOTE: ここでclarifaiのアプリケーションを使うために必要なAPI Keyの情報をctxに入れる。
	ctx := metadata.AppendToOutgoingContext(
		context.Background(),
		"Authorization", API_KEY,
	)

	// This is a publicly available model ID.
	var GeneralModelId = "aaa03c23b3724a16a56b629203edc62c"

	// NOTE: L28で作ったAPI Keyの情報(ctx)を第一引数に入れて、第二引数には検証したいモデルのIDと入力データの情報を入れる。
	response, err := client.PostModelOutputs(
		ctx,
		&api.PostModelOutputsRequest{
			ModelId: GeneralModelId,
			Inputs: []*api.Input{
				{
					Data: &api.Data{
						Image: &api.Image{
							Url: "https://samples.clarifai.com/dog2.jpeg",
						},
					},
				},
			},
		},
	)

	if err != nil {
		panic(err)
	}

	// NOTE: 返ってきたresponseで正常に検証できたかどうかを判定。できなかったらログとして出力する。
	if response.Status.Code != status.StatusCode_SUCCESS {
		panic(fmt.Sprintf("Failed response: %s", response))
	}

	// NOTE: 認識できたConseptを順にログとして表示。
	for _, concept := range response.Outputs[0].Data.Concepts {
		fmt.Printf("%s: %.2f\n", concept.Name, concept.Value)
	}
}
