import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:front/app/components/responsive.dart';
import 'package:front/constants.dart';

import 'package:get/get.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:loading_indicator/loading_indicator.dart';

import '../controllers/init_controller.dart';

class InitView extends GetView<InitController> {
  const InitView({Key? key}) : super(key: key);

  Widget renderBody(BuildContext context, {required double maxWidth}) {
    var edgeAddrss = '';
    var accessToken = '';
    var agentName = '';

    return ConstrainedBox(
      constraints: BoxConstraints(maxWidth: maxWidth),
      child: Form(
        key: controller.formKey,
        autovalidateMode: AutovalidateMode.disabled,
        child: Column(
          children: [
            const Icon(Icons.android, size: 120),
            if (!Responsive.isMobile(context)) const SizedBox(height: 10),
            Text(
              'Etri Smart Farm',
              style: GoogleFonts.bebasNeue(
                fontSize: 36,
              ),
            ),
            const SizedBox(height: 10),
            const Text(
              'Welcome our smart farm service!',
              style: TextStyle(
                fontWeight: FontWeight.bold,
                fontSize: 20,
              ),
            ),
            const SizedBox(height: 50),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 25.0),
              child: Container(
                decoration: BoxDecoration(
                  color: Colors.grey[200],
                  border: Border.all(color: Colors.white),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Padding(
                  padding: const EdgeInsets.only(left: 8.0),
                  child: TextFormField(
                    keyboardType: TextInputType.number,
                    inputFormatters: [
                      FilteringTextInputFormatter.allow(
                          RegExp(r'[0-9a-zA-Z.://]')),
                      // for version 2 and greater youcan also use this
                      // FilteringTextInputFormatter.digitsOnly
                    ],
                    decoration: const InputDecoration(
                      hintText: "Edge Address",
                      border: InputBorder.none,
                    ),
                    validator: (val) {
                      if (val!.isEmpty) {
                        return 'Edge Address는 필수사항입니다.';
                      }

                      return null;
                    },
                    onSaved: (val) {
                      edgeAddrss = val!;
                    },
                  ),
                ),
              ),
            ),
            const SizedBox(height: 10),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 25.0),
              child: Container(
                decoration: BoxDecoration(
                  color: Colors.grey[200],
                  border: Border.all(color: Colors.white),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Padding(
                  padding: const EdgeInsets.only(left: 8.0),
                  child: TextFormField(
                    keyboardType: TextInputType.number,
                    decoration: const InputDecoration(
                      hintText: "Agent Name",
                      border: InputBorder.none,
                    ),
                    validator: (val) {
                      if (val!.isEmpty) {
                        return 'Edge Address는 필수사항입니다.';
                      }
                      if (val.contains(' ')) {
                        return '공백을 포함할 수 없습니다.';
                      }
                      return null;
                    },
                    onSaved: (val) {
                      if (val != null) {
                        // controller.email = val;
                        agentName = val;
                      }
                    },
                  ),
                ),
              ),
            ),
            const SizedBox(height: 10),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 25.0),
              child: Container(
                decoration: BoxDecoration(
                  color: Colors.grey[200],
                  border: Border.all(color: Colors.white),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Padding(
                  padding: const EdgeInsets.only(left: 8.0),
                  child: TextFormField(
                    obscureText: true,
                    decoration: const InputDecoration(
                      hintText: "Access Token",
                      border: InputBorder.none,
                    ),
                    onSaved: (val) {
                      if (val != null) {
                        // controller.password = val;
                        accessToken = val;
                      }
                    },
                    validator: (val) {
                      if (val!.isEmpty) {
                        return 'Access Token은 필수사항입니다.';
                      }

                      if (val.length < 8) {
                        return '8자 이상 입력해주세요!';
                      }
                      return null;
                    },
                  ),
                ),
              ),
            ),
            const SizedBox(height: 10),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 25.0),
              child: TextButton(
                style: ButtonStyle(
                  backgroundColor: MaterialStateProperty.all(Colors.deepPurple),
                  shape: MaterialStateProperty.all<RoundedRectangleBorder>(
                    RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12)),
                  ),
                ),
                onPressed: () async {
                  if (!controller.formKey.currentState!.validate()) {
                    return;
                  }
                  controller.formKey.currentState!.save();

                  showDialog(
                    context: context,
                    builder: (BuildContext context) {
                      return const AlertDialog(
                        elevation: 0,
                        backgroundColor: Colors.transparent,
                        content: Center(
                          child: SizedBox(
                            width: 100,
                            height: 100,
                            child: LoadingIndicator(
                              indicatorType:
                                  Indicator.ballTrianglePathColoredFilled,
                              colors: kDefaultRainbowColors,
                            ),
                          ),
                        ),
                      );
                    },
                  );
                  var result = await controller.initUpdate(
                    edgeAddress: edgeAddrss,
                    agentName: agentName,
                    accessToken: accessToken,
                  );
                  Future.delayed(const Duration(seconds: 1), () {
                    Navigator.pop(context, true);
                    if (result == null) {
                      // Future.delayed(const Duration(milliseconds: 200), () {
                      //   Get.offAllNamed("/controller");
                      // });
                      Future.delayed(const Duration(milliseconds: 800), () {
                        Get.snackbar("Success", "Agent information is updated");
                      });
                      return;
                    }

                    Get.snackbar("Failed", result);
                  });
                },
                child: const SizedBox(
                  width: double.infinity,
                  child: Padding(
                    padding: EdgeInsets.all(10.0),
                    child: Center(
                      child: Text(
                        "Update",
                        style: TextStyle(
                          color: Colors.white,
                          fontWeight: FontWeight.bold,
                          fontSize: 18,
                        ),
                      ),
                    ),
                  ),
                ),
              ),
            ),
            const SizedBox(height: 25),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const Text(
                  'Do you want to init this agent?',
                  style: TextStyle(
                    fontWeight: FontWeight.bold,
                  ),
                ),
                InkWell(
                  customBorder: const StadiumBorder(),
                  child: const Padding(
                    padding:
                        EdgeInsets.symmetric(horizontal: 10.0, vertical: 20.0),
                    child: Text(
                      ' Init Now',
                      style: TextStyle(
                        color: Colors.blue,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                  onTap: () async {
                    if (!controller.formKey.currentState!.validate()) return;
                    controller.formKey.currentState!.save();

                    showDialog(
                      context: context,
                      builder: (BuildContext context) {
                        return const AlertDialog(
                          elevation: 0,
                          backgroundColor: Colors.transparent,
                          content: Center(
                            child: SizedBox(
                              width: 100,
                              height: 100,
                              child: LoadingIndicator(
                                indicatorType:
                                    Indicator.ballTrianglePathColoredFilled,
                                colors: kDefaultRainbowColors,
                              ),
                            ),
                          ),
                        );
                      },
                    );
                    var result = await controller.init(
                      edgeAddress: edgeAddrss,
                      agentName: agentName,
                      accessToken: accessToken,
                    );
                    Future.delayed(const Duration(seconds: 1), () {
                      Navigator.pop(context, true);
                      if (result == null) {
                        // Future.delayed(const Duration(milliseconds: 200), () {
                        //   Get.offAllNamed("/controller");
                        // });
                        Future.delayed(const Duration(milliseconds: 800), () {
                          Get.snackbar(
                              "Success", "Agent information is initialized");
                        });
                        return;
                      }

                      Get.snackbar("Failed", result);
                    });
                  },
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      resizeToAvoidBottomInset: true,
      body: LayoutBuilder(
        builder: (ctx, ctis) {
          return Material(
            color: Colors.transparent,
            child: SafeArea(
              child: Center(
                child: SingleChildScrollView(
                  child: Container(
                    child: ctis.maxWidth >= 700
                        ? renderBody(context, maxWidth: 700)
                        : SingleChildScrollView(
                            scrollDirection: Axis.horizontal,
                            child: renderBody(context, maxWidth: ctis.maxWidth),
                          ),
                  ),
                ),
              ),
            ),
          );
        },
      ),
    );
  }
}
