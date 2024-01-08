# Train Ticket Booking using the go and grpc  

A Small project to communicate between the client and server.

Initialy i created the directories called server, client and proto.

create the mod file in the Root directory using the below link.

go mod init github.com/HARISH-code/traindemo/ticket

Firstly i created the ticket.proto in the pb directory then i run the below command to generate the pb file.
protoc --go_out=. --go-grpc_out=.proto/ticket.proto

Then i created the main.go files for server and the client directories.

And run the code in the command prompt using  "go run main.go" commands for both server and the client.



![image](https://github.com/HARISHA-code/traindemo/assets/70417383/d4216b4e-d94a-46c5-b8e4-11423f601d1f)

