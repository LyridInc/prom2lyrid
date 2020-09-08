.\WiX\heat.exe dir ..\web\build -nologo -cg AppFiles0 -gg -g1 -srd -sfrag -template fragment -dr APPDIR0 -var var.SourceDir0 -out AppFiles0.wxs
.\WiX\heat.exe dir ..\docs -nologo -cg AppFiles1 -gg -g1 -srd -sfrag -template fragment -dr APPDIR1 -var var.SourceDir1 -out AppFiles1.wxs
.\WiX\candle.exe -dSourceDir0=..\web\build -dSourceDir1=..\docs AppFiles0.wxs AppFiles1.wxs LicenseAgreementDlg_HK.wxs WixUI_HK.wxs product.wxs
.\WiX\light.exe -ext WixUIExtension -ext WixUtilExtension -sacl -spdb  -out ..\wix\prom2lyrid.msi AppFiles0.wixobj AppFiles1.wixobj LicenseAgreementDlg_HK.wixobj WixUI_HK.wixobj product.wixobj
