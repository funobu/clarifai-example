package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Clarifai/clarifai-go-grpc/proto/clarifai/api"
	"github.com/Clarifai/clarifai-go-grpc/proto/clarifai/api/status"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func main() {
	// NOTE: カレントディレクトリの.envを環境変数とする。
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	// NOTE: crarifaiを使うのに必要な情報を環境変数からアクセスする。(カレントディレクトリに.envを作ってそこに記載する。)
	// example: API_KEY=charizard
	var API_KEY = os.Getenv("API_KEY")
	var MODEL_ID = os.Getenv("MODEL_ID")

	// NOTE: ここでclarifaiのAPIに対してgRPC通信を用いてコネクションを繋いでいる。
	conn, err := grpc.Dial(
		"api.clarifai.com:443",
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
	)
	if err != nil {
		log.Fatalf("ERROR: コネクションに失敗 (原因) -> %+v", err)
	}
	client := api.NewV2Client(conn)

	// NOTE: ここでclarifaiのアプリケーションを使うために必要なAPI Keyの情報をctxに入れる。
	ctx := metadata.AppendToOutgoingContext(
		context.Background(),
		"Authorization", "Key "+API_KEY,
	)

	// NOTE: 入力データのURL
	inputUrl := "https://pbs.twimg.com/media/FQKrdbwaMAARkpl?format=jpg&name=large"

	// NOTE: L28で作ったAPI Keyの情報(ctx)を第一引数に入れて、第二引数には検証したいモデルのIDと入力データの情報を入れる。
	response, err := client.PostModelOutputs(
		ctx,
		&api.PostModelOutputsRequest{
			ModelId: MODEL_ID,
			Inputs: []*api.Input{
				{
					Data: &api.Data{
						Image: &api.Image{
							Url: inputUrl,
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
		log.Fatalf("ERROR: 異常なレスポンスを検知 (原因) -> %s", response)
	}

	// NOTE: 認識できたConseptを順にログとして表示。
	for _, concept := range response.Outputs[0].Data.Concepts {
		fmt.Printf("%s: %.2f\n", concept.Name, concept.Value)
	}
}
