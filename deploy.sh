echo "Ready to deploy?"

echo "Stopping the service..."

ssh -i ~/coding/aws/yogabuntu.pem ubuntu@blitzapp.co "sudo service blitz stop"

echo "Stopped the service!"

echo "Recompilling the code..."
go build -o api
echo "Recompilled the code!"

echo "Uploading the new binary..."
scp -i ~/coding/aws/yogabuntu.pem api ubuntu@blitzapp.co:~/

echo "Uploaded the new binary!"

echo "Restarting the service..."
ssh -i ~/coding/aws/yogabuntu.pem ubuntu@blitzapp.co "sudo service blitz start"

echo "Restarted the service!"