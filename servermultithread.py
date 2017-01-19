import socket
import sys
import os
import threading

def read_string_from_socket(connection):
    try:
        data = connection.recv(4096)
    except socket.error, e:
        print e
        sys.exit(1)
    else:
        if len(data) == 0:
            print "Closed connection from the client"
            sys.exit(0)
        else:
            #print >>sys.stderr, 'received %s' % data
            return data

def write_string_to_socket(connection, message):
    #print >>sys.stderr, 'sending message back to the client'
    connection.sendall(message)
    
def process_client_connection(connection, client_number):

    messages = []
    while True:
        # read message 
        msg = read_string_from_socket(connection)
        message = msg.rstrip()
        print "Client %s Message received = %s " % (client_number, message)
        #sys.stdout.flush()
        messages.append(msg)
        #print eval('message == "END"')
        if message == "END":
            #print "MATCHED END"
            for msg in messages[:-1]:
                write_string_to_socket(connection, msg)
            connection.close()
            break


server_address = './uds_socket'

try:
    os.unlink(server_address)
except OSError:
    if os.path.exists(server_address):
        raise

sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
print >>sys.stderr, 'Starting up on Unix Domain Socket %s' % server_address
print >>sys.stderr, 'Waiting for connections  e.g. use socat -d - UNIX-CONNECT:./uds_socket. . .'

sock.bind(server_address)

sock.listen(1)

client_threads = []
client_number = 0

while True:
    connection, client_address = sock.accept()
    try:
        client_number = client_number + 1
        print >>sys.stderr, 'Connection from Client %s ' % client_number
        t = threading.Thread(target=process_client_connection, args=(connection,client_number,))
        client_threads.append(t)
        print >>sys.stderr, 'Starting Thread for Client %s' % client_number
        t.start()
    finally:
        if connection:
            connection.close
        
