# frozen_string_literal: true

class ApiClient
  attr_reader :base_url, :response_body, :response_status

  def initialize
    @base_url = ENV.fetch('API_METADATA_URL')
  end

  def create_client_org(npi)
    binding.pry
  end
end
