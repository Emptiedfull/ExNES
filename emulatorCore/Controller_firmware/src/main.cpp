#include <Arduino.h>
#include <Wire.h>
#include <Adafruit_GFX.h>
#include <Adafruit_SSD1306.h>

#define WIDTH 128
#define HEIGHT 32

const int rowPin[2] = {13,25};
const int colPins[3] = {12, 27, 26};

const char* buttonMap[2][3] = {
    {"A","B","UP"},{"DOWN","LEFT","RIGHT"}
};

Adafruit_SSD1306 display(WIDTH, HEIGHT, &Wire, -1);

struct Button {
    char* action;
    bool buttonPressed;
};

void setup()
{
    Serial.begin(115200);
    Wire.begin(5, 4);

    Serial.print("begginign display");

    if (!display.begin(SSD1306_SWITCHCAPVCC, 0x3C))
    {
        Serial.print("something wrong");
        for (;;)
            ;
    };
    display.clearDisplay();
    display.setTextSize(1);
    display.setTextColor(WHITE);
    display.println("display on");
    display.display();

    for (int i = 0; i < 2; i++)
    {
        pinMode(rowPin[i], OUTPUT);
        digitalWrite(rowPin[i], HIGH);
    };

    for (int i = 0; i < 3; i++)
    {
        pinMode(colPins[i], INPUT_PULLUP);
    };
};

void loop()
{
    for (int i = 0; i < 2; i++)
    {
        digitalWrite(rowPin[i], LOW);

        for (int j = 0; j < 3; j++)
        {
            if (digitalRead(colPins[j]) == LOW)
            {
                display.clearDisplay();
                display.setCursor(0, 0);

                display.print("button pressed:");
                display.println(buttonMap[i][j]);
                display.display();
            };
        };

        digitalWrite(rowPin[i], HIGH);
    };
    delay(10);
};