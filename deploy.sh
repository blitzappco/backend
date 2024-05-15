echo "Ready to deploy?"

echo "Stopping the service..."

ssh -i ~/coding/aws/yogabuntu.pem ubuntu@backend.blitzapp.co "sudo service blitz stop"

echo "Stopped the service!"

echo "Recompilling the code..."
go build -o backend
echo "Recompilled the code!"

echo "Uploading the new binary..."
scp -i ~/coding/aws/yogabuntu.pem backend ubuntu@backend.blitzapp.co:~/

echo "Uploaded the new binary!"

echo "Restarting the service..."
ssh -i ~/coding/aws/yogabuntu.pem ubuntu@backend.blitzapp.co "sudo service blitz start"

echo "Restarted the service!"