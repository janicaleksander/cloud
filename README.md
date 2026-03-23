## Car Damage Insurance Claim System 

- **Claim Service (Entry Point)**
    - **Role:** The entry point for the client to fill out the claim form and attach files (photos of the damage).
    - **Action:** Receives the request, logs the uploaded file's name and extension (no further file processing needed).
    - **Database:** Saves the new claim to its own database (e.g., `DB_Claims`) with an initial status of `NEW`.
    - **Broker:** Publishes a `ClaimSubmitted` event.

- **Policy Verification Service**
    - **Role:** Checks if the user has a valid insurance policy for the damaged car.
    - **Action:** Reacts to the `ClaimSubmitted` event.
    - **Database:** Reads from its own database (`DB_Policies`) containing information about users and their active insurances.
    - **Broker:** Publishes a `PolicyVerified` or `PolicyDenied` event based on the verification result.

- **Valuation Service**
    - **Role:** Calculates the amount of money the insurance company has to pay for the reported damages.
    - **Action:** Reacts to the `PolicyVerified` event.
    - **Database:** Reads pricing data from its own database (`DB_Valuations` - e.g., price lists for specific car parts) and saves the generated valuation record.
    - **Broker:** Publishes a `ValuationCalculated` event containing the calculated amount.

- **Decision Service**
    - **Role:** Makes the final (e.g.automated) decision regarding the payment for the damage.
    - **Action:** Reacts to the `ValuationCalculated` event. Checks if the amount is within acceptable limits.
    - **Database:** Saves the final decision and payout details into its own database (`DB_Decisions`).
    - **Broker:** Publishes a `PayoutApproved` or `PayoutRejected` event.

- **Notification Service**
    - **Role:** Handles communication with the client.
    - **Action:** Listens to all key events from the broker (`ClaimSubmitted`, `PolicyVerified`, `PolicyDenied`,`PayoutApproved`,`PayoutRejected`).
    - **Database:** Saves a history of all sent messages into its own database (`DB_Notifications`).
    - **Output:** Simulates sending information (e.g., to the user's email) by writing the message content to the application logger.




==========================
Na podstawie BPMN — zestawienie mikroserwisów i ich eventów:

---

**Claim Service**
- publikuje: `ClaimSubmitted`
- subskrybuje: `PolicyVerified`, `PolicyDenied`, `PayoutApproved`, `PayoutRejected` (zmiana statusu w DB_Claims)

---

**Policy Verification Service**
- subskrybuje: `ClaimSubmitted`
- publikuje: `PolicyVerified`, `PolicyDenied`

---

**Valuation Service**
- subskrybuje: `PolicyVerified`
- publikuje: `ValuationCalculated`

---

**Decision Service**
- subskrybuje: `ValuationCalculated`
- publikuje: `PayoutApproved`, `PayoutRejected`

---

**Notification Service**
- subskrybuje: `ClaimSubmitted`, `PolicyVerified`, `PolicyDenied`, `PayoutApproved`, `PayoutRejected`
- nie publikuje nic