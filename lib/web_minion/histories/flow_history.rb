require "web_minion/histories/history"

module WebMinion
  class FlowHistory < WebMinion::History
    attr_reader :action_history, :step_history

    def initialize
      super()
      @action_history = []
      @step_history = []
    end
  end
end
