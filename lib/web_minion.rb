require "web_minion/drivers/mechanize"
require "web_minion/step"
require "web_minion/action"
require "web_minion/flow"
require "web_minion/histories/history"
module WebMinion
  class MultipleOptionsFoundError < StandardError; end
  class NoInputFound < StandardError; end


  def self.parse_json(json)
    flow = JSON.parse(json)
    legacy_flow = {
      config: flow["config"],
      flow: {
        actions: []
      }
    }

    flow["actions"].each_with_index do |action, i|
      on_success = action["success"] || (flow["actions"][i + 1] || {})["key"] || "WebMinion_Finalizer"

      # TODO: add key automatically
      name = action["key"] || "#{action['type']} #{action["target"]}"
      legacy_action = {
        key: name,
        name: name,
        on_failure: action.fetch("error") { "WebMinion_FlowError" },
        on_success: on_success
      }

      # binding.pry

      legacy_action[:starting] = true if i.zero?
      steps = []

      case action["type"]
      when "navigate"
        step = {
          method: "go",
          target: action["target"],
          name: "Navigate to #{action['target']}"
        }
        steps << step

      when "form"

        target = {}
        if action["target"][0] == "#"
          target[:id] = action["target"][1..-1]
        elsif action["target"][0] == "."
          target[:class] = action["target"][1..-1]
        else
          target[:xpath] = action["target"]
        end

        step = {
          method: "get_form",
          # target: action["target"],
          target: target,
          name: "Get form #{action['target']}"
        }
        steps << step

        action["fields"].each do |input_field, value|
          target = {}
          if input_field["target"][0] == "#"
            target[:id] = input_field["target"][1..-1]
          elsif input_field["target"][0] == "."
            target[:class] = input_field["target"][1..-1]
          else
            target[:xpath] = input_field["target"]
          end

          steps << {
            name: "Fill in #{input_field['target']}",
            method: "fill_in_input",
            retain_element: true,
            target: target,
            value: value
          }
        end

        if action["click"]
          steps << {
            name: "Click login Button",
            method: "submit",
            target: action["click"]
          }
        else
          steps << {
            name: "Submit login",
            method: "submit"
          }
        end

        # to do click submit button step
      when "scrape"
        # target = {}
        # if action["target"][0] == "#"
        #   target[:id] = action["target"][1..-1]
        # elsif action["target"][0] == "."
        #   target[:class] = action["target"][1..-1]
        # else
        #   target[:xpath] = action["target"]
        # end

        step = {
          name: "save value #{action['target']}",
          method: "save_value",
          target: action["target"], #target,
          value: action["value"]
        }
        steps << step

        if action["value_type"]
          steps << {
            name: "format value",
            method: "format_saved_value",
            target: "balance",
            value: {
              type: action["value_type"],
              regex: action["regex"]
            }
          }
        end

      when "error"
        # TODO
      end

      if action["sleep"]
        steps << {
          name: "Sleep",
          method: "wait",
          value: 2
        }
      end

      # Add expectation step
      expectation = action["expects"] || {}
      # TODO: web_minion needs a html_contains_element validator
      validation_step = {
        value: expectation["value"] || "",
        is_validator: true,
        method: "body_includes",
        name: "validation for #{action['key']}"
      }

      if expectation["type"] == "includes"
        validation_step["method"] = "body_includes"
        # else
        # TODO: other types of validations
        # NOOP
      end

      steps << validation_step
      legacy_action[:steps] = steps
      legacy_flow[:flow][:actions] << legacy_action
    end

    # add WebMinion_Finalizer action
    # add WebMinion_FlowError action

    legacy_flow[:flow][:actions] << {
      name: "WebMinion_Finalizer",
      key: "WebMinion_Finalizer",
      steps: [
        {
          name: "WebMinion_Finalizer",
          method: "save_page_html",
          is_validator: true,
          value: "#{flow['name']}_SUCCESSFULL-INSERT_DATE.html"
        }
      ]
    }

    legacy_flow[:flow][:actions] << {
      name: "WebMinion_FlowError",
      key: "WebMinion_FlowError",
      steps: [
        {
          name: "WebMinion_FlowError",
          method: "save_page_html",
          is_validator: true,
          value: "#{flow['name']}_FAILURE-INSERT_DATE.html"
        }
      ]
    }

    legacy_flow.to_json
  end
end
