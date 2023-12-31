/*
    Video: https://www.youtube.com/watch?v=oCMOYS71NIU
    Based on Neil Kolban example for IDF:
   https://github.com/nkolban/esp32-snippets/blob/master/cpp_utils/tests/BLE%20Tests/SampleNotify.cpp
    Ported to Arduino ESP32 by Evandro Copercini

   Create a BLE server that, once we receive a connection, will send periodic
   notifications. The service advertises itself as:
   6E400001-B5A3-F393-E0A9-E50E24DCCA9E Has a characteristic of:
   6E400002-B5A3-F393-E0A9-E50E24DCCA9E - used for receiving data with "WRITE"
   Has a characteristic of: 6E400003-B5A3-F393-E0A9-E50E24DCCA9E - used to send
   data with  "NOTIFY"

   The design of creating the BLE server is:
   1. Create a BLE Server
   2. Create a BLE Service
   3. Create a BLE Characteristic on the Service
   4. Create a BLE Descriptor on the characteristic
   5. Start the service.
   6. Start advertising.

   In this example rxValue is the data received (only accessible inside that
   function). And txValue is the data to be sent, in this example just a byte
   incremented every second.
*/
#include <BLE2902.h>
#include <BLEDevice.h>
#include <BLEServer.h>
#include <BLEUtils.h>

BLECharacteristic *pCharacteristic;
bool respondedTime = 0;
bool deviceConnected = false;
float txValue = 0;
const int readPin = 32; // Use GPIO number. See ESP32 board pinouts
const int LED = 2; // Could be different depending on the dev board. I used the
                   // DOIT ESP32 dev board.

// std::string rxValue; // Could also make this a global var to access it in
// loop()

// See the following for generating UUIDs:
// https://www.uuidgenerator.net/

#define SERVICE_UUID "0000fff0-0000-1000-8000-00805f9b34fb" // UART service UUID
#define SERVICE_UUID_2 "00001801-0000-1000-8000-00805f9b34fb"
#define CHARACTERISTIC_UUID_TX "0000fff1-0000-1000-8000-00805f9b34fb"

#define CHARACTERISTIC_UUID_RX "0000fff2-0000-1000-8000-00805f9b34fb"

class MyServerCallbacks : public BLEServerCallbacks {
  void onConnect(BLEServer *pServer) {
    Serial.println("Connected!!");
    deviceConnected = true;
  };

  void onDisconnect(BLEServer *pServer) {
    Serial.println("disconnected!!");
    deviceConnected = false;
    ESP.restart();
  }
};

class MyCallbacks : public BLECharacteristicCallbacks {
  void onWrite(BLECharacteristic *rCharacteristic) {
    std::string rxValue = rCharacteristic->getValue();

    if (rxValue.length() > 0) {
      Serial.println("*********");
      Serial.print("Received Value: ");

      for (int i = 0; i < rxValue.length(); i++) {
        Serial.print(rxValue[i], HEX);
        Serial.print(" ");
      }

      Serial.println();
      Serial.println(rxValue.length());

      uint8_t txString[] = {0xFF, 0xFF, 0xFD, 0x00, 0x08,
                            0x04, 0x06, 0x0A}; // make sure this is big enuffz
      pCharacteristic->setValue(txString, 8);
      pCharacteristic->notify(); // Send the value to the app!
      respondedTime = millis();

      Serial.println("*********");
    }
  }

  // void onRead(BLECharacteristic* pCharacteristic){

  // }
};

void setup() {
  Serial.begin(115200);

  pinMode(LED, OUTPUT);

  // Create the BLE Device
  BLEDevice::init("spac1083"); // Give it a name

  // Create the BLE Server
  BLEServer *pServer = BLEDevice::createServer();
  pServer->setCallbacks(new MyServerCallbacks());

  // Create the BLE Service
  BLEService *pService = pServer->createService(SERVICE_UUID_2);
  pService = pServer->createService(SERVICE_UUID);

  // Create a BLE Characteristic
  pCharacteristic = pService->createCharacteristic(
      CHARACTERISTIC_UUID_TX, BLECharacteristic::PROPERTY_NOTIFY);

  pCharacteristic->addDescriptor(new BLE2902());

  BLECharacteristic *pCharacteristic = pService->createCharacteristic(
      CHARACTERISTIC_UUID_RX, BLECharacteristic::PROPERTY_WRITE);

  pCharacteristic->setCallbacks(new MyCallbacks());

  // Start the service
  pService->start();

  // Start advertising
  pServer->getAdvertising()->addServiceUUID(SERVICE_UUID);
  pServer->getAdvertising()->addServiceUUID(SERVICE_UUID_2);
  pServer->getAdvertising()->start();
  Serial.println("Waiting a client connection to notify...");
}

void loop() {
  int i;
  if (deviceConnected && millis() - respondedTime > 2000) {
    // Fabricate some arbitrary junk for now...
    uint8_t txString[] = {0xFF, 0xFF, 0xFD, 0x00, 0x1C, 0x23, 0x00, 0x64,
                          0x44, 0,    0,    0,    0,    0,    0,    0,
                          0,    0,    0,    0,    0,    0,    0,    0,
                          0,    0,    0,    203}; // make sure this is big
                                                  // enuffz
    pCharacteristic->setValue(txString, 28);
    pCharacteristic->notify(); // Send the value to the app!
  }
  delay(1000);
}
