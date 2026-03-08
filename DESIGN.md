# WebMinion Schema Design Document

## 1. Overview

WebMinion is a metadata-driven browser automation library that allows users to define browser automation flows through JSON configuration files instead of writing code directly. The system provides a declarative approach to creating web bots by defining flows, actions, steps, and validators.

## 2. Schema Structure

### 2.1 Configuration (config)

The config object defines the runtime environment for the automation flow:

```json
{
  "config": {
    "driver": "mechanize",
    "dimensions": {
      "width": 1024,
      "height": 768
    },
    "timeout": 30,
    "user_agent": "Mozilla/5.0"
  }
}
```

- `driver`: Specifies the browser driver to use ("mechanize" or "capybara")
- `dimensions`: Browser window dimensions (for Capybara)
- `timeout`: Request timeout in seconds
- `user_agent`: Custom user agent string

### 2.2 Flow Structure (flow)

The flow contains the main automation logic:

```json
{
  "flow": {
    "name": "Login Flow",
    "actions": [
      // Array of action objects
    ]
  }
}
```

## 3. Actions

Actions represent a group of steps that perform an automation task with validation.

### 3.1 Action Properties

```json
{
  "name": "Login Action",
  "key": "login_action",
  "starting": true,
  "on_success": "success_action",
  "on_failure": "error_action",
  "steps": [
    // Array of step objects
  ]
}
```

- `name`: Human-readable name for the action
- `key`: Unique identifier for the action
- `starting`: Boolean indicating if this action should start the flow
- `on_success`: Key of next action on success
- `on_failure`: Key of next action on failure
- `steps`: Array of step objects to execute

### 3.2 Action Types

Actions can be:
- **Starting**: The first action to execute in a flow
- **Intermediate**: Actions that link to other actions
- **Terminal**: Actions with no "on_success" field (end of execution)

## 4. Steps

Steps represent individual operations that the bot performs in an action.

### 4.1 Step Properties

```json
{
  "name": "Navigate to page",
  "method": "go",
  "target": "https://example.com/login",
  "value": "some_value",
  "retain_element": true,
  "is_validator": false,
  "skippable": false
}
```

- `name`: Human-readable description of the step
- `method`: The operation to perform (e.g., "go", "fill_in_input", "submit")
- `target`: Target element identifier or URL
- `value`: Value to input or compare against
- `retain_element`: Whether to pass the resulting element to subsequent steps
- `is_validator`: Whether this step is a validation check
- `skippable`: Whether to skip step if variables are missing

### 4.2 Valid Methods

#### Main Methods:
- `go`: Navigate to URL
- `get_form`: Get form handle by target
- `fill_in_input`: Fill input field with value
- `submit`: Submit a form
- `click`: Click an element
- `click_button_in_form`: Click button inside form
- `select`: Select from dropdown options
- `set_file_upload`: Upload a file
- `url_equals`: Validate current URL equals value
- `body_includes`: Check if page body includes text
- `value_equals`: Compare element value with expected value
- `save_page_html`: Save HTML to file
- `save_value`: Save element's value for later use
- `wait`: Wait for specified time
- `format_saved_value`: Format saved value (e.g., integer, float)

#### Select methods:
- `select/radio_button`: Select radio button
- `select/checkbox`: Check checkbox
- `select/first_radio_button`: Select first available radio button

## 5. Validators

Validators are special steps that verify automation results and determine success/failure.

### 5.1 Validator Properties

```json
{
  "name": "Verify login successful",
  "is_validator": true,
  "method": "body_includes",
  "value": "Welcome Back"
}
```

- Must have `is_validator: true` set
- Usually the last step in an action
- Return boolean status to determine flow continuation

### 5.2 Validator Methods:
- `body_includes`: Check text inclusion in page body
- `url_equals`: Compare current URL with expected value
- `value_equals`: Compare element value with expected value

## 6. Drivers

WebMinion supports multiple browser automation drivers:

### 6.1 Mechanize Driver
- Fast, lightweight, no JavaScript support
- Uses Mechanize gem for HTTP requests and HTML parsing
- Best for simple, static websites

### 6.2 Capybara Driver
- Full browser automation with JavaScript support
- Supports Selenium and Poltergeist drivers
- Can handle dynamic content and modern web applications

## 7. Variables System

Variables allow dynamic configuration:

```json
{
  "vars": {
    "username": "user@example.com",
    "password": "secret123"
  }
}
```

### 7.1 Variable Usage in Steps:
```json
{
  "value": "@username",
  "target": {
    "id": "email"
  }
}
```

- Variables are referenced with `@variable_name` syntax
- Values are substituted at runtime

## 8. Architecture

### 8.1 Core Components

**Flow**: The top-level container that coordinates execution across actions
**Action**: Groups steps and defines flow control (success/failure jumps)
**Step**: Individual bot operations (navigating, filling forms, etc.)
**Bot**: Driver abstraction that executes actual browser commands
**Validator**: Specialized steps for verification

### 8.2 Data Flow

1. Flow reads JSON configuration
2. Flow builds Action objects from configured steps
3. Flow validates flow graph for cycles and entry points
4. Flow begins execution at starting action
5. Each action executes its steps sequentially
6. Validators determine success/failure status
7. Flow navigates to appropriate next action based on result

### 8.3 Error Handling

- **Step-level errors**: Steps raise exceptions that result in failure
- **Variable missing**: Skippable steps can be skipped if variables are undefined
- **Flow validation errors**: Cyclic flows and missing starting actions are rejected before execution
- **Execution errors**: Caught and result in false status return

## 9. Testing Approach

### 9.1 Test Coverage Areas
- Flow creation from JSON configuration
- Step execution with various drivers
- Action flow graph validation
- Validator success/failure determination
- Variable substitution
- Error condition handling

### 9.2 Test Scenarios
1. Basic flow with single action
2. Multi-action flow with proper routing
3. Validation step failure
4. Variable substitution in steps
5. Driver-specific step execution
6. Error handling for missing dependencies

## 10. Implementation Plan

### Phase 1: Core Components
- Implement Flow class with JSON parsing
- Develop Action and Step class structures
- Create validator system 
- Build bot driver interfaces

### Phase 2: Validation Features
- Add cycle detection in flow graphs
- Implement starting action validation
- Build step validation system (validators must be final steps)
- Add proper error reporting

### Phase 3: Driver Support
- Complete Mechanize driver implementation
- Complete Capybara driver implementation
- Add driver-specific configuration options

### Phase 4: Advanced Features
- Variable substitution system
- Error handling improvements
- Logging and history tracking
- Extension points for custom steps

## 11. Configuration Examples

### Minimal Configuration
```json
{
  "config": {
    "driver": "mechanize"
  },
  "flow": {
    "name": "Test Flow",
    "actions": [
      {
        "key": "start",
        "starting": true,
        "steps": [
          {
            "method": "go",
            "target": "https://example.com"
          }
        ]
      }
    ]
  }
}
```

### Full Flow Example
```json
{
  "config": {
    "driver": "mechanize"
  },
  "flow": {
    "name": "Login Flow",
    "actions": [
      {
        "name": "Load Login Page",
        "key": "load_login",
        "starting": true,
        "steps": [
          {
            "name": "Go to login page",
            "method": "go",
            "target": "https://example.com/login"
          },
          {
            "name": "Get form",
            "method": "get_form",
            "target": {
              "name": "login"
            }
          },
          {
            "name": "Fill username",
            "method": "fill_in_input",
            "retain_element": true,
            "target": { "id": "email" },
            "value": "@username"
          },
          {
            "name": "Fill password",
            "method": "fill_in_input",
            "retain_element": true,
            "target": { "id": "password" },
            "value": "@password"
          },
          {
            "name": "Submit form",
            "method": "submit"
          },
          {
            "name": "Verify authenticated",
            "is_validator": true,
            "method": "body_includes",
            "value": "Welcome Back"
          }
        ]
      }
    ]
  }
}
```