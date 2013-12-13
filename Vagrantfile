# vim:filetype=ruby

Vagrant.configure('2') do |config|
  config.vm.hostname = 'hustle'
  config.vm.box = 'precise64'
  config.vm.box_url = 'http://cloud-images.ubuntu.com/vagrant/precise/' <<
                      'current/precise-server-cloudimg-amd64-vagrant-disk1.box'

  config.vm.network :private_network, ip: '33.33.33.10', auto_correct: true
  %w(8661 8662 8663 8664 8665).map(&:to_i).each do |port|
    config.vm.network :forwarded_port, guest: port, host: port,
                                       auto_correct: true
  end

  config.vm.synced_folder '.', '/gopath/src/github.com/joshk/hustle'
  config.vm.provision :shell, path: '.vagrant-provision.sh'
end
