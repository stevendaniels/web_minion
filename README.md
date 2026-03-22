# WebMinion

WebMinion is a metadata-driven browser automation library. Rather than writing a custom bot in code, you describe a browser automation flow in a JSON configuration file and hand it to WebMinion to execute. Flows are composed of named actions and steps — each step targets a page element, performs a method (navigate, fill, click, validate), and passes control to the next action.

## Implementations

| Language | Description | Directory |
|---|---|---|
| Ruby | Gem supporting Mechanize and Capybara/Selenium drivers | [ruby/webminion](ruby/webminion) |
| Go | net/http + Chrome (via go-rod) | [go-webminion](go-webminion) |

## Spec Example

A minimal flow that navigates to a URL and validates the resulting page:

```json
{
  "config": {
    "driver": "mechanize"
  },
  "flow": {
    "name": "Login Check",
    "actions": [
      {
        "name": "Go to login page",
        "starting": true,
        "steps": [
          {
            "name": "Navigate",
            "target": "https://example.com/login",
            "method": "go",
            "is_validator": false
          },
          {
            "name": "Confirm login page loaded",
            "target": "//h1[text()='Sign in']",
            "method": "element_exists",
            "is_validator": true
          }
        ]
      }
    ]
  }
}
```

See [DESIGN.md](DESIGN.md) for the full schema reference.

## License

[MIT](LICENSE.txt) · [Code of Conduct](CODE_OF_CONDUCT.md)
