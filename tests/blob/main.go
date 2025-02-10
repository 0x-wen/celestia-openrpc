package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"sync"

	openrpc "github.com/celestiaorg/celestia-openrpc"
	"github.com/celestiaorg/celestia-openrpc/types/blob"
	"github.com/celestiaorg/celestia-openrpc/types/share"
)

func Single() {
	ctx := context.Background()
	client, err := openrpc.NewClient(ctx, "ws://192.168.0.196:20002",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJwdWJsaWMiLCJyZWFkIiwid3JpdGUiLCJhZG1pbiJdfQ.OG5ptA9d6Dn4V4qRiRttX0_yUDHWk5wqZtoy18b2lIM")
	if err != nil {
		log.Fatal(err)
	}
	namespace, err := share.NewBlobNamespaceV0([]byte{1, 2, 3, 4, 5, 6, 7, 8})

	data, data2 := []byte("hello world"), []byte("hello world2")
	blobBlob, err := blob.NewBlobV0(namespace, data)
	blobBlob2, err := blob.NewBlobV0(namespace, data2)
	// blob.commitment
	fmt.Println("Blob0.Commitment:", hex.EncodeToString(blobBlob.Commitment))
	fmt.Println("Blob1.Commitment:", hex.EncodeToString(blobBlob2.Commitment))

	// write blob to DA
	height, err := client.Blob.Submit(ctx, []*blob.Blob{blobBlob, blobBlob2}, blob.NewSubmitOptions())
	fmt.Println("height: ", height)

	// retrieve data back from DA
	daBlob, err := client.Blob.Get(ctx, height, namespace, blobBlob.Commitment)

	// get all
	blobs, err := client.Blob.GetAll(ctx, height, []share.Namespace{namespace})
	for i, v := range blobs {
		fmt.Println("Blob:", i, "-", hex.EncodeToString(v.Commitment))
	}

	// get proof
	proof, err := client.Blob.GetProof(ctx, height, namespace, daBlob.Commitment)

	// Included
	included, err := client.Blob.Included(ctx, height, namespace, proof, daBlob.Commitment)
	fmt.Println("Included0:", included)

	proof2, err := client.Blob.GetProof(ctx, height, namespace, daBlob.Commitment)
	included2, err := client.Blob.Included(ctx, height, namespace, proof2, daBlob.Commitment)
	fmt.Println("Included1:", included2)

	// get commitment proof
	commitmentProof, err := client.Blob.GetCommitmentProof(ctx, height, namespace, daBlob.Commitment)
	fmt.Println("CommitmentProof0:", commitmentProof)
	commitmentProof2, err := client.Blob.GetCommitmentProof(ctx, height, namespace, blobBlob2.Commitment)
	fmt.Println("CommitmentProof1:", commitmentProof2)
}

func getAll() {
	ctx := context.Background()
	client1, err := openrpc.NewClient(ctx, "ws://192.168.0.196:20002",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJwdWJsaWMiLCJyZWFkIiwid3JpdGUiLCJhZG1pbiJdfQ.OG5ptA9d6Dn4V4qRiRttX0_yUDHWk5wqZtoy18b2lIM")
	if err != nil {
		log.Fatal(err)
	}
	client2, err := openrpc.NewClient(ctx, "ws://192.168.0.196:20012",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJwdWJsaWMiLCJyZWFkIiwid3JpdGUiXX0.MJ0Aj2WXicZDrt9-BqsAqpovjor6RpByhzPFEMGsqp8")
	if err != nil {
		log.Fatal(err)
	}

	namespace, err := share.NewBlobNamespaceV0([]byte{1, 2, 3, 4, 5, 6, 7, 8})

	data, data2 := []byte("hello world"), []byte("hello world2")
	blobBlob, err := blob.NewBlobV0(namespace, data)
	blobBlob2, err := blob.NewBlobV0(namespace, data2)

	wg := sync.WaitGroup{}
	wg.Add(2)

	ch := make(chan uint64, 2)
	// write blob to DA
	go func() {
		defer wg.Done()
		height1, err := client1.Blob.Submit(ctx, []*blob.Blob{blobBlob}, blob.NewSubmitOptions())
		if err != nil {
			log.Fatal(err)
		}
		ch <- height1
	}()
	submitOpt := blob.NewSubmitOptions(blob.WithGasPrice(0.002))
	fmt.Println("GasPrice: ", submitOpt.GasPrice())
	go func() {
		defer wg.Done()
		height2, err := client2.Blob.Submit(ctx, []*blob.Blob{blobBlob2}, submitOpt)
		if err != nil {
			log.Fatal(err)
		}
		ch <- height2
	}()

	wg.Wait()
	close(ch)

	var heights []uint64
	for i := range ch {
		heights = append(heights, i)
	}
	if len(heights) != 2 {
		log.Fatal("Expected exactly two heights, but got ", len(heights))
	}
	if !allElementsEqual(heights) {
		log.Fatal("Expected all heights to be equal, but they are not")
	}

	// Use one of the heights for GetAll
	daBlob, err := client1.Blob.GetAll(ctx, heights[0], []share.Namespace{namespace})
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range daBlob {
		matchData := bytes.Equal(data, v.Data)
		matchCommitment1 := bytes.Equal(v.Commitment, blobBlob.Commitment)
		matchCommitment2 := bytes.Equal(v.Commitment, blobBlob2.Commitment)

		// 先检查数据是否符合要求
		if !(matchData && matchCommitment1) && !(!matchData && matchCommitment2) {
			log.Fatal("data is error")
		}

		// 统一获取 proof
		proof, err := client1.Blob.GetProof(ctx, heights[0], namespace, v.Commitment)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Proof:", proof)

		// 统一检查 Included 状态
		commitment := blobBlob.Commitment
		if !matchData {
			commitment = blobBlob2.Commitment
		}

		included, err := client1.Blob.Included(ctx, heights[0], namespace, proof, commitment)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Included:", included)

	}

	fmt.Println("Blob:", daBlob)
}

func allElementsEqual(slice []uint64) bool {
	if len(slice) == 0 {
		return true
	}

	first := slice[0]
	for _, v := range slice[1:] {
		if v != first {
			return false
		}
	}
	return true
}

func main() {
	Single()
	//getAll()
}
