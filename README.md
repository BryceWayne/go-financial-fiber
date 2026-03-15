# Go Financial Dashboard

A sleek, dynamic personal finance dashboard built with Go, Fiber, and HTMX.

## Overview

This application provides a comprehensive, single-page financial calculator that breaks down your pre-tax income, applies estimated tax calculations based on your location (State/City), and categorizes your monthly expenses and savings. It features a modern, responsive UI with a glassmorphism aesthetic.

## Features

- **Real-time Calculations**: Powered by HTMX, the dashboard updates your financial breakdown instantly as you type, without requiring a full page reload.
- **Accurate Tax Estimates**: Includes marginal tax bracket calculations for Federal taxes, with specific state and city tax logic (e.g., NY, CA, IL, NYC, SF).
- **Categorized Finances**:
  - **Income & Taxes**: Pre-tax annual/monthly, effective tax rates, and post-tax take-home pay.
  - **Core Expenses**: Mortgage, groceries, utilities, student loans, etc.
  - **Discretionary Spending**: Dining out, travel, subscriptions, etc.
  - **Leisure & Savings**: Fun money, 401k, cash savings, brokerage, and home equity.
- **Bonus Handling**: Separates bonus income and applies appropriate tax rates so you can see your true bonus take-home pay.
- **Modern UI**: Custom CSS utilizing CSS variables, flexbox/grid layouts, and a dark "glass" theme with dynamic background animations.

## Tech Stack

- **Backend**: [Go](https://go.dev/) (1.22+)
- **Web Framework**: [Fiber v2](https://gofiber.io/)
- **Templating**: Go `html/template` (via Fiber template engine)
- **Frontend Interactivity**: [HTMX](https://htmx.org/)
- **Styling**: Vanilla CSS

## Getting Started

### Prerequisites

- Go 1.22.2 or higher installed on your machine.

### Installation

1. Clone the repository and navigate to the project directory:
   ```bash
   cd go-financial-fiber
   ```

2. Install the necessary Go dependencies:
   ```bash
   go mod tidy
   ```

3. Run the application:
   ```bash
   go run main.go
   ```

4. Open your web browser and navigate to:
   ```
   http://localhost:3000
   ```

## Project Structure

- `main.go`: The core backend logic. Contains the Fiber server setup, tax bracket definitions, the calculation logic (`calculateTotals`), and route handlers.
- `views/`: Contains the HTML templates for the frontend.
  - `index.html`: The main layout, CSS styling, and HTMX inclusion.
  - `dashboard.html`: The interactive form and data grid that HTMX re-renders on input changes.
- `go.mod` / `go.sum`: Go module dependencies tracking.

## Developing

To modify the tax brackets or add new locations, update the `TaxBracket` structs and switch statements inside the `calculateTotals` function in `main.go`.

To tweak the UI layout or colors, modify the CSS variables or styles inside `views/index.html`.
