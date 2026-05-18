# 🚜 Smart Farm Edge-to-Cloud Data Pipeline & Control HMI

This repository demonstrates a production-grade, edge-to-cloud IoT data pipeline designed for Agritech environments. It features decoupled architecture, utilizing a lightweight edge node for sensor metrics and a centralized Kubernetes (K3s) cluster for data aggregation, monitoring, and real-time HMI control.

## 🏗️ Core Architecture (Decoupled System)

The infrastructure is built with cloud-native principles, strictly separating the edge computing node (Farm) from the centralized monitoring plane (Control Tower) to prevent Single Points of Failure (SPOF) and ensure high availability in unstable network environments.

### 1. Edge Node (`farm-node`: 192.168.202.131)
* **Role:** Responsible for local sensor data ingestion and hardware actuation.
* **Component:** Golang Mock Sensor (Generates and exposes Temperature, Humidity, and Soil Moisture metrics).
* **Design Choice:** Exposes metrics via an HTTP endpoint (`/metrics`) to allow the centralized server to scrape data (Pull Model), preventing overload on the edge device.

### 2. Control Tower (`monitor-node`: 192.168.202.132)
* **Role:** Centralized data aggregation, visualization, and custom alert routing.
* **Infrastructure:** K3s (Lightweight Kubernetes) cluster for container orchestration and self-healing.
* **Components:**
    * **Prometheus:** Scrapes time-series data from the Edge Node.
    * **Grafana (NodePort: 31165):** Visualizes metrics and triggers alerts.
    * **Golang Webhook Server:** Custom backend service that receives Grafana alerts, formats them via Go Templates, and pushes notifications to **Telegram**.
    * **Node-RED:** Serves as the Human-Machine Interface (HMI) for real-time digital twin visualization and manual override controls (e.g., Force Cooling Actuation).

## 🔄 Data Pipeline Flow
```mermaid
graph LR
    subgraph EdgeNode [Farm-Node: 192.168.202.131]
        direction TB
        Sensor[Golang Sensor Daemon<br/>Generates data every 2s]
        LocalLogic{Local Control Logic}
        Actuator[Sprinkler / Ventilation]
        
        Sensor -->|Auto-Trigger| LocalLogic
        LocalLogic -->|Local Actuation| Actuator
    end

    subgraph K3sCluster [Monitor-Node: 192.168.202.132]
        direction TB
        Prometheus[Prometheus Server]
        Grafana[Grafana Dashboard]
        NodeRED[Node-RED HMI]
        Webhook{{Golang Webhook Server}}
        
        Prometheus -->|2. Query Metrics| Grafana
        Prometheus -->|3. Query Metrics| NodeRED
        Grafana -->|4. Trigger Alert| Webhook
    end

    Sensor -->|1. HTTP Pull| Prometheus
    Webhook -->|5. Format and Push| Telegram([Telegram App])
    NodeRED -->|6. Manual Override| Actuator
