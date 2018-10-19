module WebMinion
  class Bot
    attr_reader :config
    attr_accessor :bot

    def initialize(config = {})
      @config = config
    end

    def execute_step(method, target, value = nil, element = nil, values_hash = {})
      if [:save_value, :format_saved_value].include?(method)
        method(method).call(target, value, element, values_hash)
      else
        method(method).call(target, value, element)
      end
    end
  end
end
