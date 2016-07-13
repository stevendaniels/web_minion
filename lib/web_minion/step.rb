module WebMinion
  class InvalidMethodError < StandardError; end
  class NoValueForVariableError < StandardError; end
  # A Step represents the individual operation that the bot will perform. This
  # often includes grabbing an element from the DOM tree, or performing some
  # operation on an element that has already been found.
  class Step
    attr_accessor :name, :target, :method, :value, :is_validator, :retain_element
    attr_reader :saved_values, :vars

    VALID_METHODS = {
      select: [
        :field,
        :radio_button,
        :first_radio_button,
        :checkbox
      ],
      main_methods: [
        :get_field,
        :get_form,
        :go,
        :select,
        :click,
        :click_button_in_form,
        :submit,
        :fill_in_input,
        :url_equals,
        :value_equals,
        :body_includes,
        :save_page_html,
        :save_value
      ]
    }.freeze

    def initialize(fields = {})
      fields.each_pair do |k, v|
        if valid_method?(k.to_sym)
          send("method=", k)
          @target = v
        else
          send("#{k}=", v)
        end
      end

      replace_all_variables
    end

    def vars=(vars)
      @vars = Hash[vars.collect{|k,v| [k.to_s, v]}] 
    end 

    def perform(bot, element = nil, saved_values)
      bot.execute_step(@method, @target, @value, element, saved_values)
    end

    def method=(method)
      raise(InvalidMethodError, "Method: #{method} is not valid") unless valid_method?(method.to_sym)
      split = method.to_s.split("/").map(&:to_sym)
      @method = split.count > 1 ? VALID_METHODS[split[0]][split[1]] : method.to_sym
    end

    def retain?
      retain_element
    end

    def validator?
      is_validator
    end

    def valid_method?(method)
      split = method.to_s.split("/").map(&:to_sym)
      if split.count > 1
        return true unless VALID_METHODS[split[0]][split[1]].nil?
      end
      VALID_METHODS[:main_methods].include?(method)
    end
  
    private

    def replace_all_variables
      %w(value target).each do |field|
        return if send(field).nil?
        if replace_var = send(field).match(/@(\D+)/)
          raise(NoValueForVariableError, "no variable to use found for #{replace_var}") unless @vars[replace_var[1]]
          send("#{field}=", @vars[replace_var[1]])
        end
      end  
    end
  end
end
