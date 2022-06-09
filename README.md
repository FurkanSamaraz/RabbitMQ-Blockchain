# RabbitMQ-Blockchain

RabbitMQ install

1- docker

2- docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.10-management

We send our data from the client to rabbitmq and read the sent data from the server that is already listening. We obtain the blocked data by inserting the read data into the blockchain-based encryption algorithm.

İstemciden verimizi rabbitmq' e yolluyoruz ve yollanan veriyi hazırdan dinleyen serverden okuyoruz. Okunan veriyi blockchain tabanı şifreleme algoritmasına sokarak block'lanmış veriyi elde ediyoruz.

![RabbitMQ](https://user-images.githubusercontent.com/92402372/172896360-d9a0272f-d161-49ee-9f9b-40508d82c6aa.png)


<img width="571" alt="Ekran Resmi 2022-06-09 19 27 51" src="https://user-images.githubusercontent.com/92402372/172897511-d6191976-1aae-465e-aadf-e59e6ba9f04e.png">
