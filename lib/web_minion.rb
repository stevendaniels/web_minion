require "web_minion/drivers/mechanize"
require "web_minion/step"
require "web_minion/action"
require "web_minion/flow"
require "web_minion/histories/history"
module WebMinion
  class MultipleOptionsFoundError < StandardError; end
  class NoInputFound < StandardError; end
end
