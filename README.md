# hpc-torque-helper
A TCP-socket server and client-side tools for retrieving Torque/Moab job information that requires elevated privileges, such as root or the Torque system admin.

## How it works
The server is designed to run on the Torque server.  It should be run by an user who has the required privileges to retrieve job information from the torque/moab server.

The server listen on a TCP socket, and waiting for a client to send in a *command* following a commnication protocol (see below). When the server recieves a command, it maps the command to a system call and returns output to the client.

## cliet-server communication protocol
1. Client initiates a TCP socket with TLS.
1. Client sends a command string to the server in format of `<cmdName>++++<cmdArg1>++++<cmdArg2>` followed by an ending character `\n` (ASCII control character 0010) to indicate the end of the command.
1. Server receives the command string until the `\n` character.
1. Server performs the command (mapped to the system call) on the server side.
1. Server sends the output to the client, and the character `\a` (ASCII control character 0007) to indicate the end of the output.
1. Client receives the output until the `\a` character.
1. Client repeats with the next command or sends the command `bye\n` to close the connection.
