package grpcapi

import (
	"context"
	"database/sql"
	"log"
	"net"
	"time"

	"mailinglist-ms/maildb"
	pb "mailinglist-ms/proto"

	"google.golang.org/grpc"
)

type MailServer struct {
	pb.UnimplementedMailingListServiceServer
	db *sql.DB
}

// conver protbuf message

func pbEntryToMailDbEntry(pbEntry *pb.EmailEntry) maildb.EmailEntry{
	tme := time.Unix(pbEntry.ConfirmedAt,0)
	return maildb.EmailEntry{
		Id: pbEntry.Id,
		Email: pbEntry.Email,
		ConfirmedAt: &tme,
		OptOut: pbEntry.OptOut,
	}

}


func mailDbToPbEntry(maildbEntry *maildb.EmailEntry)pb.EmailEntry{
	return pb.EmailEntry{
		Id: maildbEntry.Id,
		Email: maildbEntry.Email,
		ConfirmedAt: maildbEntry.ConfirmedAt.Unix(),
		OptOut: maildbEntry.OptOut,
	}
}




func emailResponse(db *sql.DB,email string) (*pb.EmailResponse, error){
	entry , err := maildb.GetEmail(db,email)
	if err != nil {
		return &pb.EmailResponse{} ,err 
	}
	if entry == nil {
		return &pb.EmailResponse{}, nil
	}
	
	response := mailDbToPbEntry(entry);

	return &pb.EmailResponse{EmailEntry: &response},nil
}



func (s *MailServer) GetEmail(ctx context.Context, req *pb.GetEmailRequest)(*pb.EmailResponse,error){
	log.Printf("grpc GetEmail: %v\n",req)
	return emailResponse(s.db,req.EmailAddr)
}

func (s *MailServer) GetEmailBatch(ctx context.Context, req *pb.GetEmailBatchRequest)(*pb.GetEmailBatchResponse,error){
	log.Printf("grpc GetEmailBatch:%v\n",req)
	 
	params := maildb.GetEmailBatchQueryParameters{
		Page: int(req.Page),
		Count: int(req.Count),
	}

	
	maildbEntries, err := maildb.GetEmailBatch(s.db,params)
	if err != nil{
		return &pb.GetEmailBatchResponse{},err
	}

	pbEntries := make([]*pb.EmailEntry,0,len(maildbEntries))
	for i:=range maildbEntries{
		entry := mailDbToPbEntry(&maildbEntries[i])
		pbEntries = append(pbEntries, &entry)
	}

		return &pb.GetEmailBatchResponse{EmailEntries: pbEntries}, nil
}


func (s *MailServer) CreateEmail(ctx context.Context, req *pb.CreateEmailRequest)(*pb.EmailResponse,error){
	log.Printf("grpc CreateEmail: %v\n",req)

	err := maildb.CreateEmail(s.db,req.EmailAddr)
	if err != nil {
		return &pb.EmailResponse{}, err
	}
	return emailResponse(s.db,req.EmailAddr)
}

func (s *MailServer) UpdateEmail(ctx context.Context, req *pb.UpdateEmailRequest)(*pb.EmailResponse,error){
	log.Printf("grpc UpdateEmail: %v\n",req)
	entry := pbEntryToMailDbEntry(req.EmailEntry)
	err := maildb.UpdateEmail(s.db,entry)
	if err != nil {
		return &pb.EmailResponse{}, err
	}
	return emailResponse(s.db,entry.Email)
}

func (s *MailServer) DeleteEmail(ctx context.Context, req *pb.DeleteEmailRequest)(*pb.EmailResponse,error){
	log.Printf("grpc DeleteEmail: %v\n",req)
	err := maildb.DeleteEmail(s.db,req.EmailAddr)
	if err != nil {
		return &pb.EmailResponse{}, err
	}
	return emailResponse(s.db,req.EmailAddr)
}

func Serve(db *sql.DB,bind string){
	listener, err := net.Listen("tcp",bind)
	if err != nil {
		log.Fatalf("gRPC server error: failed to bind %v\n",bind)
	}

	grpcServer := grpc.NewServer()
	mailServer := MailServer{db:db}
	pb.RegisterMailingListServiceServer(grpcServer, &mailServer)


	log.Printf("gRPC API server listening on %v\n",bind)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("gRPC server error: %v\\n",err)
	}


}
