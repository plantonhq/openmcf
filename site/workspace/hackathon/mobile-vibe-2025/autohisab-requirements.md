# AutoHisab - Mobile App Requirements

**App Name:** AutoHisab (Auto = Vehicle, Hisab = Calculation/Accounting in Hindi)

**Document Purpose:** This is a comprehensive requirements document for building the AutoHisab mobile app. All product decisions, technical choices, and implementation guidelines are captured here for AI-assisted development.

---

## 1. App Overview

### What is AutoHisab?

AutoHisab is a fuel refill tracker that calculates vehicle mileage (fuel efficiency). Every time you refuel, you log the odometer reading and fuel quantity. The app calculates distance traveled and mileage between refills, showing fuel efficiency trends over time.

### Core Value Proposition

- Simple fuel tracking without complexity
- Automatic mileage calculation between refills
- Visual trends to understand fuel efficiency over time
- Multiple vehicle support for personal use

### Target User

**Primary User Type:** Personal use - individual tracking their own vehicles

**User Profile:**
- Single user (no family sharing or collaboration)
- May track multiple vehicles (5+ vehicles supported)
- Refuels weekly (typical usage pattern)
- Wants to understand fuel efficiency trends
- May want to track fuel expenses

---

## 2. Product Requirements

### MVP Features (Must-Have for Phase 1)

**Essential features that define the minimum viable product:**

1. **Log Fuel Refill**
   - Enter odometer reading (keyboard input)
   - Enter fuel quantity (keyboard input)
   - Capture date/time automatically
   - Optional: fuel type (petrol, diesel, CNG)
   - Optional: notes field

2. **Auto-Calculate Mileage**
   - Distance = Current odometer - Previous odometer
   - Mileage = Distance / Fuel quantity
   - Display immediately after logging refill

3. **View Refill History**
   - List of all refills (newest first)
   - Each entry shows: date, odometer, fuel quantity, calculated mileage
   - Organized by vehicle

4. **View Mileage Trends**
   - Line chart showing mileage efficiency over time
   - Time period: last 6 months or all data
   - Visual identification of improving/worsening trends

5. **Multiple Vehicle Support**
   - Add/manage multiple vehicles
   - Track refills per vehicle
   - Switch between vehicles easily
   - Compare mileage across vehicles

6. **Analytics Dashboard**
   - Current mileage (latest calculated value)
   - Average mileage (over time period)
   - Best mileage (highest efficiency recorded)
   - Worst mileage (lowest efficiency recorded)
   - Total fuel expenses (this month/year)
   - Cost per kilometer (if cost data available)

### Explicitly OUT OF SCOPE for Phase 1

**These features are intentionally excluded to keep the app simple:**

- ❌ Backend/API integration
- ❌ User authentication or accounts
- ❌ Multi-device sync
- ❌ Data backup or export (CSV, PDF)
- ❌ Camera access / Photo of odometer or receipt
- ❌ OCR (Optical Character Recognition)
- ❌ Voice input
- ❌ GPS/location tracking
- ❌ Gas station name or location
- ❌ Background processing
- ❌ Reminders or notifications
- ❌ Maintenance schedule tracking
- ❌ Partial refill handling (assume full tank)
- ❌ City vs Highway mileage comparison
- ❌ Summer vs Winter comparison (though date is captured)
- ❌ App store deployment
- ❌ Automated testing

---

## 3. User Workflows

### Primary User Flow: Logging a Refill

**Scenario:** User just refueled at a gas station, still in the car

**Quick Entry Flow:**
1. User opens AutoHisab
2. User sees Dashboard with prominent "Add Refill" button
3. User taps "Add Refill"
4. Form displays with fields:
   - Vehicle selector (if multiple vehicles)
   - Odometer reading (number keyboard)
   - Fuel quantity (number keyboard with decimal)
   - Fuel type (dropdown: Petrol/Diesel/CNG) - optional
   - Date/time (auto-populated, editable if needed)
   - Notes (text field) - optional
5. User enters odometer reading
6. User enters fuel quantity
7. User taps "Save"
8. App calculates mileage immediately
9. App displays: "Your mileage: X.X km/L"
10. User returns to dashboard with updated stats

**Expected Time:** Under 30 seconds for experienced user

### Home Screen: Dashboard

**What user sees when app opens:**

- Header with app name and selected vehicle
- Large card showing:
  - Latest mileage (big, prominent number)
  - "X.X km/L" or "X.X MPG"
- Quick stats section:
  - Average mileage (last 10 refills)
  - Total expense (this month)
- Recent refills list (last 3-5 entries):
  - Date
  - Odometer reading
  - Fuel quantity
  - Calculated mileage
- Prominent "Add Refill" button (floating action button)
- Navigation to other screens:
  - Analytics/Charts
  - Refill History (full list)
  - Vehicles Management
  - Settings

### Secondary Flows

**View Mileage Trends:**
1. From dashboard, tap "Analytics" or "View Chart"
2. See line graph of mileage over time
3. Time period selector (1 month, 3 months, 6 months, 1 year, all)
4. Summary stats below chart

**Manage Vehicles:**
1. From dashboard, tap vehicle selector or settings
2. See list of all vehicles
3. Add new vehicle: name, make, model, year
4. Edit or delete existing vehicles
5. Set active vehicle for quick refill logging

**View Full Refill History:**
1. From dashboard, tap "View All" or "History"
2. See chronological list of all refills
3. Filter by vehicle
4. Tap entry to see details or edit

---

## 4. Data Model

### Core Entities

#### Vehicle

```
Vehicle {
  id: string (UUID)
  name: string (e.g., "Honda City", "My Bike")
  make: string (e.g., "Honda")
  model: string (e.g., "City")
  year: number (optional)
  odometerUnit: string ("km" or "miles")
  fuelUnit: string ("liters" or "gallons")
  createdAt: DateTime
}
```

#### Refill Entry

```
RefillEntry {
  id: string (UUID)
  vehicleId: string (foreign key to Vehicle)
  odometerReading: number (in km or miles)
  fuelQuantity: number (in liters or gallons)
  fuelCost: number (optional, total amount spent)
  pricePerUnit: number (optional, calculated or manual)
  fuelType: string (optional: "petrol", "diesel", "cng")
  dateTime: DateTime
  notes: string (optional)
  calculatedDistance: number (computed)
  calculatedMileage: number (computed)
  createdAt: DateTime
}
```

### Mileage Calculation

**Standard calculation formula:**

```
Distance = Current Odometer Reading - Previous Odometer Reading
Mileage = Distance / Fuel Quantity

Example:
Current refill: Odometer 45,000 km, Fuel 40 liters
Previous refill: Odometer 44,500 km
→ Distance: 45,000 - 44,500 = 500 km
→ Mileage: 500 / 40 = 12.5 km/liter
```

**Edge Cases to Handle:**

1. **First refill for a vehicle:**
   - No previous reading exists
   - Cannot calculate mileage
   - Display: "First refill - baseline set"
   - Store the entry but no mileage calculation

2. **Invalid data (current < previous):**
   - Show error: "Odometer reading must be greater than previous reading"
   - Don't allow save

3. **Zero or negative fuel quantity:**
   - Show error: "Fuel quantity must be greater than zero"
   - Don't allow save

4. **Odometer rollover (999,999 → 0):**
   - Out of scope for Phase 1 (rare edge case)

### Units

**Default Units:**
- Distance: Kilometers (km)
- Fuel: Liters (L)
- Mileage: km/L

**Future consideration:** Support for miles and gallons (MPG)

---

## 5. Analytics & Reports

### Core Analytics (Dashboard)

**Single number stats:**
- **Current Mileage:** Latest calculated mileage from most recent refill
- **Average Mileage:** Mean of last 10 refills (or all refills if < 10)
- **Best Mileage:** Maximum efficiency ever recorded for this vehicle
- **Worst Mileage:** Minimum efficiency recorded for this vehicle

### Financial Analytics

**If fuel cost is captured:**
- **Total Fuel Expenses:**
  - This month
  - This year
  - All-time
- **Cost per Kilometer:**
  - Total fuel cost / Total distance traveled
- **Fuel Expense Trends:**
  - Monthly bar chart of total expenses

### Comparison Analytics

**Across Vehicles:**
- If user has multiple vehicles, show comparison:
  - Vehicle A: 15.2 km/L average
  - Vehicle B: 18.5 km/L average
  - Vehicle C: 12.1 km/L average

### Mileage Trend Chart

**Visualization Requirements:**
- **Chart Type:** Line graph
- **X-axis:** Date/time of refill
- **Y-axis:** Calculated mileage (km/L)
- **Data points:** Each refill (except first refill with no mileage)
- **Time period selector:** 1M, 3M, 6M, 1Y, All
- **Trend indicator:** Visual cue if improving/worsening/stable

---

## 6. Technical Architecture

### Platform Decision

**Platform:** Flutter

**Rationale from discussion:**
- Single codebase for iOS and Android
- Fast performance for forms and charts
- Beautiful UI out of box
- Strong hot reload for development
- Team will rely on AI (Cursor) for Dart learning curve

### Storage Strategy

**Storage:** Hive (local NoSQL database)

**Key Decisions:**
- **Local-only storage** (no backend, no cloud sync)
- **Data persists** on device using Hive
- **Data is lost** if app is deleted (acceptable for Phase 1)
- **No backup or export** for Phase 1
- **Offline-first by design** (actually offline-only)

**Why Hive:**
- Fast NoSQL database for Flutter
- Simple key-value and object storage
- No backend required
- Lightweight and performant

### Authentication & Accounts

**Authentication:** None

**Key Decisions:**
- No user accounts required
- Anonymous usage (single device)
- No login or sign-up flow
- App works immediately on first launch

### State Management

**Decision:** Let AI (Cursor) choose based on Flutter best practices

**Likely options:**
- Provider (recommended by Flutter team)
- Riverpod (modern Provider)
- Bloc (if more structure needed)

**State Requirements:**
- UI state (form inputs, modal open/closed, selected tab)
- App state (current vehicle selection, user preferences)
- Persistent state (refill entries, vehicles) → stored in Hive

### No Backend

**Confirmed Decisions:**
- No REST API
- No GraphQL
- No Firebase/Supabase
- No cloud sync
- All data local to device

---

## 7. UI/UX Requirements

### Design Philosophy

- **Simple and clean:** Minimal UI, no clutter
- **Fast data entry:** Optimized for quick refill logging
- **Clear visualization:** Mileage trends easy to understand
- **Material Design:** Follow Flutter Material guidelines

### Home Screen Layout

**Dashboard Structure:**

```
┌─────────────────────────────┐
│  AutoHisab        [Vehicle▼]│
│                             │
│  ┌───────────────────────┐ │
│  │  Latest Mileage       │ │
│  │      15.2 km/L        │ │
│  │  ↑ +0.5 from previous │ │
│  └───────────────────────┘ │
│                             │
│  Quick Stats                │
│  ┌──────────┐ ┌──────────┐ │
│  │ Avg: 14.8│ │ Best: 16.2│ │
│  └──────────┘ └──────────┘ │
│                             │
│  Recent Refills             │
│  ┌─────────────────────────┐│
│  │ Dec 20 | 45000 | 40L    ││
│  │ Mileage: 15.2 km/L      ││
│  ├─────────────────────────┤│
│  │ Dec 15 | 44500 | 38L    ││
│  │ Mileage: 14.7 km/L      ││
│  └─────────────────────────┘│
│                             │
│  [Analytics] [History] [⚙] │
│                             │
│              [+] Add Refill │
└─────────────────────────────┘
```

### Add Refill Form

**Form Fields (in order):**

1. **Vehicle Selector** (if multiple vehicles)
   - Dropdown or modal picker
   - Shows vehicle name (e.g., "Honda City")

2. **Odometer Reading** (required)
   - Label: "Odometer Reading"
   - Input type: Number keyboard
   - Placeholder: "45000"
   - Hint: Previous reading shown below (e.g., "Last: 44500 km")

3. **Fuel Quantity** (required)
   - Label: "Fuel Quantity (Liters)"
   - Input type: Decimal number keyboard
   - Placeholder: "40.5"

4. **Fuel Cost** (optional)
   - Label: "Total Cost (optional)"
   - Input type: Decimal number keyboard
   - Placeholder: "3500"

5. **Fuel Type** (optional)
   - Label: "Fuel Type"
   - Input: Segmented control or dropdown
   - Options: Petrol, Diesel, CNG

6. **Date & Time** (auto-filled, editable)
   - Default: Current date/time
   - User can tap to change if logging old refill

7. **Notes** (optional)
   - Label: "Notes (optional)"
   - Input: Text field (multiline)
   - Placeholder: "Highway trip, AC on..."

**Buttons:**
- **Save** (primary button, prominent)
- **Cancel** (secondary button)

**After Save:**
- Show success message with calculated mileage
- "Mileage: 12.5 km/L ✓"
- Return to dashboard with updated data

### Refill History Screen

**List View:**
- Chronological order (newest first)
- Group by vehicle (if multiple)
- Each entry shows:
  - Date (formatted: "Dec 20, 2024")
  - Odometer reading
  - Fuel quantity
  - Calculated mileage (prominent)
  - Fuel cost (if available)
- Tap entry to see full details or edit
- Filter option by vehicle

### Analytics/Chart Screen

**Mileage Trend Chart:**
- Line graph at top (takes up 50% of screen)
- Time period selector below chart (chips: 1M, 3M, 6M, 1Y, All)
- Summary stats below:
  - Average mileage for selected period
  - Best mileage in period
  - Worst mileage in period
  - Total distance traveled
  - Total fuel consumed

**Chart Interaction:**
- Tap data point to see details
- Pinch to zoom (optional)
- Smooth animations

### Vehicles Management Screen

**List of Vehicles:**
- Card for each vehicle showing:
  - Vehicle name
  - Make and model
  - Total refills logged
  - Average mileage
  - Last refill date
- Add new vehicle button
- Edit/delete options

---

## 8. Development Environment & Tools

### IDE

**IDE:** Cursor (AI-powered VS Code fork)
- Already installed and configured
- Primary development tool

### Platform Tools

**Flutter SDK:**
- Already set up on development machines
- Version: Latest stable

**Simulators/Emulators:**
- iOS Simulator (Mac users)
- Android Emulator (all platforms)
- Both already configured

### Device Testing

**Physical Devices:**
- Team has iOS and Android devices
- Connection setup already understood
- Testing on real devices available

### Version Control

**Git Repository:**
- **Organization:** Left Wing Mobile Apps (GitHub)
- **Repository Structure:** One repo per app
- **Repository Name:** `autohisab-app` (or similar)
- **Branching:** Simple feature branches
- **Commits:** Follow conventional commit format

---

## 9. Implementation Philosophy

### AI-First Development

**Core Principle:** Rely heavily on Cursor AI for all implementation decisions not explicitly specified in this document.

**What this means:**
- Cursor AI should choose best practices for Flutter architecture
- Cursor AI should select appropriate packages for charts, forms, navigation
- Cursor AI should implement UI following Material Design guidelines
- Cursor AI should handle edge cases reasonably
- Team will not write code manually; all code generated by AI

### Keep It Simple

**Guiding Philosophy:** Avoid complexity at all costs.

**Decision-Making Framework:**
- If a feature seems complex → defer to Phase 2
- If multiple approaches exist → choose the simplest
- If technical decision is unclear → ask AI to use Flutter best practices
- If edge case is rare → handle gracefully with error message, don't over-engineer

### Manual Testing Only

**Testing Strategy for Phase 1:**
- No automated unit tests
- No automated integration tests
- No E2E test automation
- Manual testing in emulators/simulators
- Manual testing on physical devices
- Test checklist (provided below)

### No Deployment for Phase 1

**Phase 1 Goal:** App runs successfully in local emulator/simulator

**Not in scope:**
- App store submission
- TestFlight distribution
- APK/IPA distribution
- Code signing for production
- Build optimization

---

## 10. Success Criteria for Phase 1

**The app is considered successful when:**

1. ✅ App launches successfully in iOS and Android emulators
2. ✅ User can add a new vehicle (name, make, model)
3. ✅ User can select active vehicle (if multiple exist)
4. ✅ User can log first refill (odometer, fuel quantity)
5. ✅ First refill is stored successfully (no mileage calculated, expected)
6. ✅ User can log second refill
7. ✅ Mileage is calculated correctly (distance/fuel)
8. ✅ Calculated mileage is displayed immediately after save
9. ✅ Dashboard shows latest mileage prominently
10. ✅ Dashboard shows quick stats (average, best, worst)
11. ✅ Recent refills list appears on dashboard
12. ✅ User can view full refill history
13. ✅ Refill history is sorted newest first
14. ✅ User can view mileage trend chart
15. ✅ Chart displays line graph with data points
16. ✅ All data persists using Hive (survives app restart)
17. ✅ Form validation works (odometer > previous, fuel > 0)
18. ✅ Error messages display for invalid input
19. ✅ App handles edge cases gracefully (first refill, no data)
20. ✅ UI is clean, usable, and follows Material Design

---

## 11. Manual Testing Checklist

**Before considering Phase 1 complete, manually verify:**

### Fresh Install
- [ ] App installs successfully on iOS emulator
- [ ] App installs successfully on Android emulator
- [ ] App launches without errors
- [ ] Empty state displays correctly (no vehicles, no refills)

### Vehicle Management
- [ ] Can add first vehicle (name, make, model)
- [ ] Vehicle appears in list/selector
- [ ] Can add second vehicle
- [ ] Can switch between vehicles
- [ ] Can edit vehicle details
- [ ] Can delete vehicle (with confirmation)

### Logging Refills
- [ ] "Add Refill" button is prominent and accessible
- [ ] Form displays with all required fields
- [ ] Odometer field accepts numbers only
- [ ] Fuel quantity field accepts decimals
- [ ] Date/time is auto-populated with current time
- [ ] Can select fuel type (petrol/diesel/CNG)
- [ ] Can add optional notes
- [ ] Can save refill successfully

### First Refill Handling
- [ ] First refill for vehicle saves successfully
- [ ] No mileage calculation attempted (expected)
- [ ] Appropriate message shown (e.g., "Baseline set")
- [ ] Entry appears in refill history

### Mileage Calculation
- [ ] Second refill calculates mileage correctly
- [ ] Formula verified: (current - previous) / fuel
- [ ] Calculated mileage displays immediately
- [ ] Mileage appears in refill entry
- [ ] Dashboard updates with latest mileage

### Validation
- [ ] Odometer < previous shows error (doesn't save)
- [ ] Odometer = previous shows error
- [ ] Fuel quantity = 0 shows error
- [ ] Fuel quantity < 0 shows error
- [ ] Required fields empty shows error
- [ ] Error messages are clear and helpful

### Dashboard
- [ ] Latest mileage displays prominently
- [ ] Average mileage calculates correctly
- [ ] Best mileage shows highest value
- [ ] Worst mileage shows lowest value
- [ ] Recent refills list shows last 3-5 entries
- [ ] Quick stats are accurate
- [ ] All data updates after new refill

### Refill History
- [ ] All refills display in list
- [ ] Sorted newest first
- [ ] Each entry shows date, odometer, fuel, mileage
- [ ] Can tap entry to view details
- [ ] First refill shows "no mileage" or similar
- [ ] Empty state handled (no refills yet)

### Analytics/Chart
- [ ] Mileage chart displays line graph
- [ ] Data points match refill entries
- [ ] X-axis shows dates correctly
- [ ] Y-axis shows mileage values correctly
- [ ] Chart is readable and clear
- [ ] Time period selector works (if implemented)
- [ ] Summary stats below chart are accurate

### Data Persistence
- [ ] Close app and reopen
- [ ] All vehicles still present
- [ ] All refills still present
- [ ] Dashboard data unchanged
- [ ] No data loss on app restart

### Multiple Vehicles
- [ ] Can track refills for vehicle A
- [ ] Switch to vehicle B
- [ ] Log refills for vehicle B
- [ ] Switch back to vehicle A
- [ ] Vehicle A data unchanged
- [ ] Comparison works (if implemented)

### Edge Cases
- [ ] Large odometer values (100,000+) work
- [ ] Small fuel quantities (5L) work
- [ ] Large fuel quantities (80L) work
- [ ] Decimal fuel quantities (38.5L) work
- [ ] Very first app launch handled gracefully
- [ ] No vehicles → prompts to add vehicle
- [ ] No refills → prompts to add refill

### UI/UX
- [ ] Text is readable (good contrast)
- [ ] Buttons are tappable (good size)
- [ ] Forms are easy to fill out
- [ ] Navigation is intuitive
- [ ] Back buttons work correctly
- [ ] Modals close properly
- [ ] Loading states display (if applicable)
- [ ] Success messages clear
- [ ] Error messages clear

### Cross-Platform
- [ ] Same behavior on iOS emulator
- [ ] Same behavior on Android emulator
- [ ] UI adapts appropriately to each platform
- [ ] No platform-specific crashes

---

## 12. Implementation Guidelines for AI

### What AI Should Decide

**Cursor AI has full autonomy to decide:**

- Flutter project structure (feature-based, type-based, etc.)
- State management solution (Provider, Riverpod, Bloc, etc.)
- Navigation approach (Navigator 2.0, go_router, etc.)
- Form management (built-in Form widgets, flutter_form_builder, etc.)
- Chart library (fl_chart, charts_flutter, syncfusion, etc.)
- Date formatting library (intl, jiffy, etc.)
- UI component approach (Material widgets, custom widgets)
- Code organization and architecture patterns
- Error handling patterns
- Validation logic implementation
- Any Flutter best practices not specified here

### What is Already Decided

**Cursor AI must follow these decisions:**

- **Platform:** Flutter (not React Native, not Native)
- **Storage:** Hive (not SQLite, not Realm, not Firebase)
- **No backend:** Local storage only, no API calls
- **No authentication:** No user accounts
- **Units:** Kilometers and Liters by default
- **Home screen:** Dashboard (not refill history, not chart)
- **Input method:** Keyboard only (no camera, no voice)
- **Offline:** Fully offline, no sync required
- **Calculation:** Distance = Current - Previous, Mileage = Distance / Fuel

### Code Quality Expectations

**Code should be:**
- Clean and readable
- Well-commented where logic is complex
- Following Flutter/Dart conventions
- Handling errors gracefully
- Providing helpful error messages to users
- Performant (but don't over-optimize prematurely)

### User Experience Priorities

**Prioritize:**
1. **Fast refill entry** - the primary use case should be quick
2. **Clear feedback** - user always knows what happened (success, error)
3. **Obvious navigation** - user never lost or confused
4. **Data visibility** - important numbers prominent and clear
5. **Graceful errors** - helpful messages, never crashes

---

## 13. Data Examples

### Example Vehicle Data

```json
{
  "id": "vehicle-001",
  "name": "Honda City",
  "make": "Honda",
  "model": "City",
  "year": 2020,
  "odometerUnit": "km",
  "fuelUnit": "liters"
}
```

### Example Refill Entries

**First Refill (no mileage):**
```json
{
  "id": "refill-001",
  "vehicleId": "vehicle-001",
  "odometerReading": 44500,
  "fuelQuantity": 40,
  "fuelCost": 3520,
  "pricePerUnit": 88,
  "fuelType": "petrol",
  "dateTime": "2024-12-15T10:30:00Z",
  "notes": "",
  "calculatedDistance": null,
  "calculatedMileage": null
}
```

**Second Refill (with mileage):**
```json
{
  "id": "refill-002",
  "vehicleId": "vehicle-001",
  "odometerReading": 45000,
  "fuelQuantity": 38,
  "fuelCost": 3344,
  "pricePerUnit": 88,
  "fuelType": "petrol",
  "dateTime": "2024-12-20T09:15:00Z",
  "notes": "Highway driving, AC on",
  "calculatedDistance": 500,
  "calculatedMileage": 13.16
}
```

**Third Refill:**
```json
{
  "id": "refill-003",
  "vehicleId": "vehicle-001",
  "odometerReading": 45520,
  "fuelQuantity": 35,
  "fuelCost": 3080,
  "pricePerUnit": 88,
  "fuelType": "petrol",
  "dateTime": "2024-12-25T17:45:00Z",
  "notes": "City traffic",
  "calculatedDistance": 520,
  "calculatedMileage": 14.86
}
```

### Example Analytics Calculations

**Average Mileage:**
```
Refills with mileage: [13.16, 14.86]
Average = (13.16 + 14.86) / 2 = 14.01 km/L
```

**Best Mileage:**
```
Max of [13.16, 14.86] = 14.86 km/L
```

**Total Distance:**
```
Sum of all distances = 500 + 520 = 1020 km
```

**Total Fuel:**
```
Sum of all fuel = 40 + 38 + 35 = 113 liters
```

**Total Cost:**
```
Sum of all costs = 3520 + 3344 + 3080 = 9944 rupees
```

**Cost per km:**
```
Total cost / Total distance = 9944 / 1020 = 9.75 rupees/km
```

---

## 14. Next Steps After This Document

**How to use this document:**

1. **Provide to Cursor AI** - This document is the complete specification
2. **Start implementation** - Cursor can generate the entire Flutter project
3. **Iterate on UI** - Refine UI based on what gets generated
4. **Test manually** - Use the checklist above
5. **Capture learnings** - Document what worked/didn't work for future apps

**Development Flow:**

```
This Requirements Doc → Cursor AI → Generated Flutter Project
                                   ↓
                            Manual Testing in Emulator
                                   ↓
                            Refinement/Bug Fixes (via Cursor)
                                   ↓
                            Phase 1 Complete ✓
                                   ↓
                            Phase 2 Planning (backend, sync, export, etc.)
```

---

## 15. Questions for Cursor AI

**When implementing, if unclear:**

1. **Architecture:** Use Flutter best practices for project structure
2. **State Management:** Choose the simplest solution that works (likely Provider)
3. **UI Library:** Use standard Material widgets unless custom needed
4. **Chart Library:** Choose fl_chart (popular, well-maintained)
5. **Forms:** Use built-in Form widgets with validation
6. **Navigation:** Use go_router or Navigator 2.0 (modern approach)
7. **Date Formatting:** Use intl package (official)
8. **Hive Setup:** Follow Hive documentation for model definitions

**If edge case not specified:** Handle gracefully with error message, don't crash.

---

## Appendix: Key Quotes from Transcript

**On simplicity:**
> "We want to keep this scope as simple as possible and... keep the scope as simple as possible. That's the philosophy, because this is the first mobile app that we are trying to build."

**On complexity avoidance:**
> "Background processing not required... cursor says that background processing is complex. And our funder is keep it simple. Anything that is close to complex, we want to avoid."

**On AI-first development:**
> "Team doesn't know Dart Swift Kotlin, but we rely heavily on... we are not going to team is not going to write a single line of code and we are promoting wire coding."

**On testing:**
> "We are not worried about testing strategies. We don't want any automation. We do manual testing since we are picking these small apps, micro apps."

**On backend:**
> "The app will only work offline. There won't be any back end. So data storage will be stored locally only."

**On success criteria:**
> "Right now we want this app to be built and just want to start it in one of the emulators. So that will be a good... goal that we can achieve for this hackathon."

---

**End of Requirements Document**

This document captures all decisions made during the planning discussion. It is ready to be provided to Cursor AI for implementation of the AutoHisab mobile app.



