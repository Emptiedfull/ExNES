#include <Arduino.h>

const int PIR_PIN = 13;     
const int ONBOARD_LED = 2;   

void setup() {
  Serial.begin(115200);
  
  pinMode(PIR_PIN, INPUT);   
  pinMode(ONBOARD_LED, OUTPUT);
  
  Serial.println("PIR Motion Sensor Warming Up...");
  delay(10000);              
  Serial.println("Sensor Active!");
}

void loop() {
  int motionDetected = digitalRead(PIR_PIN);

  if (motionDetected == HIGH) {
    Serial.println("MOTION DETECTED!");
    digitalWrite(ONBOARD_LED, HIGH); 
  } else {  
    digitalWrite(ONBOARD_LED, LOW);
  }
  
  delay(200); 
}