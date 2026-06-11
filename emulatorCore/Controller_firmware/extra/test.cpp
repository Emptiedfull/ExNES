

#include <Arduino.h>
#include <WiFi.h>

#define BUTTON_B 15
#define BUTTON_A 12
#define LED 25

const char* ssid = "Bikash_Goyal";
const char* password = "mishti12345";
const char* host = "192.168.29.102";
const int port = 8090;

struct Button {
  const int PIN;
  volatile bool Pressed;
  const char* action;
};

Button buttons[] = {
  {BUTTON_B, false, "BUTTON_B"},
  {BUTTON_A, false, "BUTTON_A"},
};

int pressCount = 0;


void IRAM_ATTR handleButtonChange(Button &btn){
  int state = digitalRead(btn.PIN);

  if (state == LOW){
    btn.Pressed = true;
  }else{
    btn.Pressed = false;
  }
}

void IRAM_ATTR isrB() { handleButtonChange(buttons[0]); }
void IRAM_ATTR isrA() { handleButtonChange(buttons[1]); }

WiFiClient client;

void setup() {
  Serial.begin(115200);

  

  pinMode(BUTTON_B, INPUT_PULLUP);
  pinMode(BUTTON_A, INPUT_PULLUP);
  pinMode(LED,OUTPUT);
  attachInterrupt(digitalPinToInterrupt(BUTTON_B), isrB, CHANGE);
  attachInterrupt(digitalPinToInterrupt(BUTTON_A), isrA, CHANGE);

 
  WiFi.begin(ssid,password);

  Serial.println("Connecting to wifi");

  while (WiFi.status() != WL_CONNECTED)
  {
    Serial.print(".");
    delay(50);
  };

  Serial.println("wifi connected with local IP:");
  Serial.print(WiFi.localIP());

  if (client.connect(host,port)){
    Serial.println("connected to emulator server");
    digitalWrite(LED,HIGH);
  };

  client.println("Welcome to esp32 controller");


}

unsigned long lastSent = 0;
const unsigned long interval = 5000;


void loop() {

  if (millis()-lastSent > interval){
    lastSent = millis();
    if (client.connected()){
      client.println("ping");
    }else{
      Serial.println("no connection found");
    }
  }

  
}