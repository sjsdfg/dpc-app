# frozen_string_literal: true

class ApiClient
  attr_reader :base_url

  def initialize
    @base_url = ENV.fetch('API_METADATA_URL')
  end
end
