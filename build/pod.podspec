Pod::Spec.new do |spec|
  spec.name         = 'Gfro'
  spec.version      = '{{.Version}}'
  spec.license      = { :type => 'GNU Lesser General Public License, Version 3.0' }
  spec.homepage     = 'https://github.com/frogeum/go-frogeum'
  spec.authors      = { {{range .Contributors}}
		'{{.Name}}' => '{{.Email}}',{{end}}
	}
  spec.summary      = 'iOS Frogeum Client'
  spec.source       = { :git => 'https://github.com/frogeum/go-frogeum.git', :commit => '{{.Commit}}' }

	spec.platform = :ios
  spec.ios.deployment_target  = '9.0'
	spec.ios.vendored_frameworks = 'Frameworks/Gfro.framework'

	spec.prepare_command = <<-CMD
    curl https://gfrostore.blob.core.windows.net/builds/{{.Archive}}.tar.gz | tar -xvz
    mkdir Frameworks
    mv {{.Archive}}/Gfro.framework Frameworks
    rm -rf {{.Archive}}
  CMD
end
