## Product Specification: Simple Esoteric Analysis App

---

## 1. Project Overview & Core Value Proposition
Build a minimalist, mobile-first web/native application that provides personalized Bazi, numerology, and holistic life analysis. 

*   **Core Value Proposition:** Maximum insight with minimum friction.
*   **Strategy:** **Progressive Disclosure.** The app onboards users rapidly with minimal friction and reveals deeper, high-precision insights as they provide more precise data.

---

## 2. User Journey & UI Architecture

### Phase 1: Rapid Onboarding (Step 1 Wizard)
*   **UI State:** Ultra-clean, single-input focus.
*   **Input Required:** Birthdate only ($DD/MM/YYYY$).
*   **Action:** User submits birthdate to instantly generate their base profile.

### Phase 2: The Soft Gate (Conditional Modal)
*   **Trigger:** Activated if the user attempts to move forward without entering a birth time or location.
*   **UI Component:** Sleek pop-up modal.
*   **Copy/UX Intent:** 
    > "Without your birth time, we can’t calculate your Hour Pillar (which rules your career and late-life fortune). You can still view your core personality, but career and wealth insights will be locked. Proceed anyway or add time?"
*   **Buttons:** `[ Proceed with Basic Profile ]` or `[ Add Time Now ]`

### Phase 3: The Main Dashboard Layout
The dashboard prioritizes scannability, utilizing a clear text-based hierarchy over complex spreadsheets or raw Chinese characters.
+-------------------------------------------------------------+
| [Top] Personality Summary                                   |
| "Core Profile: Balanced Wood Element / Life Path 7"         |
+-------------------------------------------------------------+
| [Middle Navigation] Horizontal Tab Bar                      |
|  💼 Career  |  ❤️ Romance  |  🍏 Health  |  💰 Wealth       |
+-------------------------------------------------------------+
| [Content Body]                                              |
| Dynamic display based on the active tab selection           |
+-------------------------------------------------------------+

---

## 3. The Visual "Lock" Mechanism
If a user bypasses Phase 2 and enters only a birthdate, the dashboard restricts advanced tabs using an interactive lock strategy.

*   **Visual Presentation:** 
    *   `[ ❤️ Romance ]` and `[ 🍏 Health ]` tabs display basic data.
    *   `[ 💼 Career ]` and `[ 💰 Wealth ]` tabs show a beautifully styled, blurred-out preview graphic of a complex Bazi chart overlayed with a secure lock icon.
*   **Call to Action (CTA):** A button reading: `[ Unlock Deep Insights with Birth Time & City ]`.
*   **UX Action:** Tapping this button launches Phase 4 (Step 2 Wizard).

---

## 4. Phase 4: Deep Analysis Onboarding (Step 2 Wizard)
*   **Inputs Required:** Birth Time ($HH:MM$) + Birth City/Country.
*   **City Input UX Requirements:** Uses a dynamic autocomplete dropdown text field to capture location data (mapping to backend latitude, longitude, and historical timezone database).
*   **Fallback Handling ("City Not Found"):** If a user types a remote location and the database returns zero results, the dropdown transforms into a helper text module:
    > “Can’t find your town? Try typing the nearest major city or regional capital instead. (This ensures your timezone calculation remains accurate!)”
*   **Final Action:** Upon submission, the dashboard refreshes seamlessly, removing all blurred elements and fully populating all tabs with high-precision calculations.

---

## 5. Design & Aesthetic Guidelines

### Style & Data Presentation
*   **Aesthetic:** Strict minimalism. Neutral backgrounds combined with high-contrast, intentional typography.
*   **Abstraction Layer:** Hide complex metaphysical math behind plain text, clear summaries, and clean visual indicators (e.g., simple elemental balance bar charts). Avoid surface-layer technical jargon.

### Micro-interactions
*   **Tab Toggling:** Smooth sliding transitions when switching between active tabs.
*   **Locked Tabs:** A subtle, glowing accent on locked tabs to prompt user interaction and engagement.