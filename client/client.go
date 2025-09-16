package main

import (
	"context"
	"log"
	pb "mailinglist-ms/proto"
	"time"

	"github.com/alexflint/go-arg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// seperate application !important



func logResponse(res *pb.EmailResponse, err error){
	if err != nil {
		log.Fatalf("error : %v",err)
	}
	if res.EmailEntry == nil {
		log.Printf("Email not found!")
	} else {
		log.Printf("Response : %v",res.EmailEntry)
	}
}


func CreateEmail(client pb.MailingListServiceClient, addr string) *pb.EmailEntry{
	log.Println("Create Email")
	ctx, cancel := context.WithTimeout(context.Background(),time.Second)
	defer cancel()
	
	res, err := client.CreateEmail(ctx, &pb.CreateEmailRequest{EmailAddr:addr})
	logResponse(res,err)

	return res.EmailEntry
}

func UpdateEmail(client pb.MailingListServiceClient, entry pb.EmailEntry) *pb.EmailEntry{
	log.Println("Update Email")
	
	ctx, cancel := context.WithTimeout(context.Background(),time.Second)
	defer cancel()
	
	res, err := client.UpdateEmail(ctx, &pb.UpdateEmailRequest{EmailEntry: &entry})
	logResponse(res,err)
	return res.EmailEntry
}

func DeleteEmail(client pb.MailingListServiceClient, addr string) *pb.EmailEntry{
	log.Println("Delete Email")
	
	ctx, cancel := context.WithTimeout(context.Background(),time.Second)
	defer cancel()
	
	res, err := client.DeleteEmail(ctx, &pb.DeleteEmailRequest{EmailAddr: addr})
	logResponse(res,err)
	return res.EmailEntry
}



func GetEmail(client pb.MailingListServiceClient, addr string) *pb.EmailEntry{
	log.Println("Get Email")
	ctx, cancel := context.WithTimeout(context.Background(),time.Second)
	defer cancel()
	
	res, err := client.GetEmail(ctx, &pb.GetEmailRequest{EmailAddr:addr})
	logResponse(res,err)
	return res.EmailEntry
}
func GetEmailBatch(client pb.MailingListServiceClient, count,page int32) {
	log.Println("Get Email Batch")
	ctx, cancel := context.WithTimeout(context.Background(),time.Second)
	defer cancel()

	res, err := client.GetEmailBatch(ctx, &pb.GetEmailBatchRequest{Page: page, Count:count})

	if err != nil {
		log.Fatalf("Error: %v",err)
	}
	log.Println("Response:")
	for i:=0 ; i< len(res.EmailEntries);i++{
		log.Printf("item [%v of %v]: %s",i+1,len(res.EmailEntries),res.EmailEntries[i])
	}

}





var args struct {
	gRPCAddr string `arg:"env:MAILINGLIST_GRPC_ADDR"`
}

func main(){
	arg.MustParse(&args)

	if args.gRPCAddr == ""{
		args.gRPCAddr = ":8082"
	}

		conn, err := grpc.Dial(args.gRPCAddr,grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Did not connect: %v",err)
	}
	defer conn.Close()


	client := pb.NewMailingListServiceClient(conn)

	newEmail := CreateEmail(client,"dandified@abc.com")
	newEmail.ConfirmedAt = 1000
	UpdateEmail(client, *newEmail)
	DeleteEmail(client,newEmail.Email)
	GetEmailBatch(client,5,1)
}
