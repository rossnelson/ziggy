# Ziggy Hardware Guide

Physical hardware planning and assembly for Ziggy devices. Workflows run in Temporal Cloud; devices are thin clients.

---

## Bill of Materials

| Part | Model | ~Cost |
|------|-------|-------|
| Single-board computer | Raspberry Pi Zero 2 W | $15 |
| Display | Pimoroni 1.3" SPI LCD (240x240, ST7789) | $15 |
| Battery SHIM | Pimoroni LiPo SHIM | $10 |
| Battery | Pimoroni LiPo 1200mAh (or similar JST) | $12 |
| Buttons | 3x 12mm tactile buttons + caps | $8 |
| SD Card | 16GB Class 10 | $8 |
| GPIO Header | 2x20 pin header (if not pre-soldered) | $2 |
| Wiring | Jumper wires | $5 |
| Case | 3D printed | $0 |
| Misc | M2.5 standoffs, screws | $5 |

**Total: ~$80 per unit**

### Why Pi Zero 2 W?

- $15 vs $45 for Pi 4
- Smaller form factor (65mm x 30mm)
- Built-in WiFi
- Sufficient for worker + display (512MB RAM)
- Temporal runs in cloud, not on device

### Battery Details

**Pimoroni LiPo SHIM:**
- Solders directly to GPIO pins (very thin, ~3mm)
- JST connector for LiPo battery
- Charges via Pi's USB port while running
- Low battery warning via GPIO 4 (optional)
- Safe shutdown support

**Battery options (JST-PH 2-pin):**
- 1200mAh → ~4-6 hours runtime
- 2000mAh → ~6-10 hours runtime
- 6000mAh → all day (bigger, needs case space)

Charge by plugging USB into Pi. Ziggy keeps running while charging.

---

## GPIO Pinout

### Display (SPI)

| LCD Pin | Pi GPIO | Pi Pin # |
|---------|---------|----------|
| VCC | 3.3V | 1 |
| GND | GND | 6 |
| SCL/SCLK | GPIO 11 | 23 |
| SDA/MOSI | GPIO 10 | 19 |
| DC | GPIO 9 | 21 |
| CS | GPIO 8 | 24 |
| BL | GPIO 13 | 33 |
| RST | GPIO 25 | 22 |

### Buttons

| Button | Function | Pi GPIO | Pi Pin # |
|--------|----------|---------|----------|
| BTN1 | Feed | GPIO 17 | 11 |
| BTN2 | Play | GPIO 27 | 13 |
| BTN3 | Pet | GPIO 22 | 15 |
| Common | Ground | GND | 9 |

---

## Wiring Diagram

```
                  Raspberry Pi Zero 2 W
                 ┌─────────────────────┐
                 │ (1) 3.3V      5V (2)│
      LCD VCC ───│                     │
                 │ (5)        GND (6) │──── LCD GND
                 │                     │
   Buttons GND ──│ (9) GND            │
    BTN1 Feed ───│(11) GPIO17         │
    BTN2 Play ───│(13) GPIO27         │
     BTN3 Pet ───│(15) GPIO22         │
                 │                     │
      LCD MOSI ──│(19) GPIO10         │
       LCD DC ───│(21) GPIO9   22(22) │──── LCD RST
      LCD CLK ───│(23) GPIO11  24(24) │──── LCD CS
                 │                     │
       LCD BL ───│(33) GPIO13         │
                 └─────────────────────┘
```

---

## WiFi Setup

### First Boot Experience

1. **Power on** → LCD shows:
   ```
   Welcome to Ziggy!
   
   Connect to WiFi:
   Ziggy-Setup
   ```

2. **User connects** phone/laptop to `Ziggy-Setup` network

3. **Captive portal opens** (or browse to `192.168.4.1`)
   - Shows available networks
   - User selects theirs, enters password

4. **Device saves config, connects to WiFi, pulls containers**

5. **Success** → LCD shows:
   ```
   Connected!
   
   View Ziggy at:
   cloud.temporal.io
   ```

### Reconfiguring WiFi

- **Option A**: Hold all 3 buttons during boot (5 sec) to force AP mode
- **Option B**: Use balenaCloud dashboard to purge WiFi config

---

## Balena Fleet

### docker-compose.yml

```yaml
version: '2'

services:
  wifi-connect:
    image: balenablocks/wifi-connect:latest
    network_mode: host
    labels:
      io.balena.features.dbus: '1'
      io.balena.features.firmware: '1'
    cap_add:
      - NET_ADMIN
    environment:
      - PORTAL_SSID=Ziggy-Setup

  ziggy:
    build: ./ziggy
    environment:
      - TEMPORAL_ADDRESS=${TEMPORAL_ADDRESS}
      - TEMPORAL_NAMESPACE=${TEMPORAL_NAMESPACE}
      - TEMPORAL_API_KEY=${TEMPORAL_API_KEY}
      - ANTHROPIC_API_KEY=${ANTHROPIC_API_KEY}
      - OWNER_NAME=${OWNER_NAME}
      - TIMEZONE=${TIMEZONE}
    ports:
      - "8080:8080"

  display:
    build: ./hardware
    depends_on:
      - ziggy
    devices:
      - "/dev/spidev0.0:/dev/spidev0.0"
      - "/dev/gpiomem:/dev/gpiomem"
    cap_add:
      - SYS_RAWIO
    environment:
      - ZIGGY_API=http://ziggy:8080

volumes: {}
```

### Device Variables

Set per-device in balenaCloud:

| Variable | Example | Purpose |
|----------|---------|---------|
| `TEMPORAL_ADDRESS` | `ziggy.abc12.tmprl.cloud:7233` | Cloud endpoint |
| `TEMPORAL_NAMESPACE` | `ziggy-prod` | Shared namespace |
| `TEMPORAL_API_KEY` | `<key>` | Auth |
| `ANTHROPIC_API_KEY` | `sk-ant-...` | AI responses |
| `OWNER_NAME` | `ross` | Workflow ID prefix |
| `TIMEZONE` | `America/Los_Angeles` | Sleep cycle |

### Fleet Variables (shared)

| Variable | Value |
|----------|-------|
| `TEMPORAL_ADDRESS` | `ziggy.abc12.tmprl.cloud:7233` |
| `TEMPORAL_NAMESPACE` | `ziggy-prod` |
| `TEMPORAL_API_KEY` | `<shared-key>` |
| `ANTHROPIC_API_KEY` | `<shared-key>` |

Only `OWNER_NAME` and `TIMEZONE` need to be per-device.

### OTA Updates

```bash
balena push ziggy-fleet
```

All devices update automatically.

---

## Temporal Cloud

### Namespace

All team Ziggys share one namespace:

```
Name: ziggy-prod
Retention: 30 days
```

### Viewing Workflows

Everyone can see all Ziggys at:

```
https://cloud.temporal.io/namespaces/ziggy-prod/workflows
```

Filter by workflow ID prefix to find your own:

```
ziggy-ross-*
```

### What Happens If Device Dies?

**Nothing.** Workflow keeps running in cloud. Scheduled decay continues. Get a new device, set same `OWNER_NAME`, it reconnects to existing workflow. Ziggy remembers everything.

This is the demo: **Durable execution means Ziggy survives hardware failure.**

---

## 3D Printed Case

### Pi Zero Form Factor

Stack order (bottom to top):
- Battery (fits in case floor)
- LiPo SHIM (soldered to Pi)
- Pi Zero 2 W
- Display (on GPIO header or wired)

Total stack height: ~25-30mm depending on battery

### Layout

```
┌──────────────────────┐
│   ┌────────────┐     │
│   │            │     │
│   │  Display   │     │
│   │            │     │
│   └────────────┘     │
│                      │
│  (🍔)  (🎮)  (💜)   │
│  Feed  Play  Pet     │
└──────────────────────┘
    ▲              ▲
    │              └── Charging LED window (optional)
    └── USB Micro-B (charging + power)
```

### Design Notes

- Battery compartment in base
- USB port accessible for charging (Pi's micro USB)
- LiPo SHIM adds minimal height (~3mm)
- Consider LED window to show charging status
- Ventilation less critical (lower power)
- Optional: on/off switch wired to SHIM

---

## Quick Start Card

```
ZIGGY SETUP

1. Plug in USB to charge/power
2. On your phone, connect to "Ziggy-Setup" WiFi
3. Enter your home WiFi password
4. Wait for Ziggy to appear!

VIEW ALL ZIGGYS
cloud.temporal.io/namespaces/ziggy-prod

BUTTONS
🍔 Feed  |  🎮 Play  |  💜 Pet

BATTERY
Unplug and take Ziggy anywhere!
Recharge anytime via USB.

RESET WIFI
Hold all 3 buttons for 5 seconds on boot

FUN FACT
Unplug your device—Ziggy keeps living in the cloud.
Plug back in, and Ziggy remembers everything!

NEED HELP?
#ziggy-support on Slack
```

---

## Assembly Checklist

- [ ] Solder GPIO header to Pi Zero (if needed)
- [ ] Solder LiPo SHIM to Pi Zero GPIO pins
- [ ] Connect display via SPI pins
- [ ] Wire buttons to GPIO 17, 27, 22 + ground
- [ ] Connect LiPo battery to SHIM (JST connector)
- [ ] Flash SD card with balenaOS
- [ ] Insert SD card, power on via USB
- [ ] Configure WiFi via captive portal
- [ ] Verify device appears in balenaCloud
- [ ] Set device variables (OWNER_NAME, TIMEZONE)
- [ ] Test buttons and display
- [ ] Unplug USB, verify battery power works
- [ ] Assemble into case (battery in base)
- [ ] Ship to team member
