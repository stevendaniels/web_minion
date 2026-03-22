require "web_minion/histories/history"

module WebMinion
  class StepHistory < WebMinion::History
    attr_reader :step_name

    def initialize(step_name)
      super()
      @step_name = step_name
    end
  end
end