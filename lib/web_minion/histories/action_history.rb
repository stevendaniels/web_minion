module WebMinion
  class ActionHistory < WebMinion::History
    attr_reader :action_name, :action_key, :step_history

    def initialize(action_name, action_key)
      super()
      @action_name = action_name
      @action_key = action_key
      @step_history = []
    end
  end
end