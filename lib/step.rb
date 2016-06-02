class InvalidMethodError < StandardError; end

class Step
  attr_accessor :name, :target, :method, :value, :is_validator, :retain_element

  VALID_METHODS = [
    :go,
    :select,
    :click,
    :get_form,
    :get_form_in_index,
    :get_field,
    :select_field,
    :select_radio_button,
    :select_checkbox,
    :submit,
    :fill_in_input,
    :url_equals,
    :value_equals,
    :body_includes,
    :save_page_html
  ].freeze

  def initialize(fields = {})
    fields.each_pair do |k, v|
      send("#{k}=", v)
    end
  end

  def perform(bot, element = nil)
    bot.execute_step(@method, @target, @value, element)
  end

  def method=(method)
    raise(InvalidMethodError, "Method: #{method} is not valid") unless valid_method?(method.to_sym)
    @method = method.to_sym
  end

  def retain?
    retain_element
  end

  def validator?
    is_validator
  end

  def valid_method?(method)
    VALID_METHODS.include?(method)
  end
end
