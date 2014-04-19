//
//  WHONewMessageViewController.h
//  Hooli
//
//  Created by dylan on 4/19/14.
//  Copyright (c) 2014 whoisdylan. All rights reserved.
//

#import <UIKit/UIKit.h>

@protocol WHOMessageProtocol <NSObject>
- (void)receivedNewMessage:(NSString* )message withLocation:(NSString* ) location;
@end

@interface WHONewMessageViewController : UIViewController
- (IBAction)messageView:(id)sender;
@property (strong, nonatomic) IBOutlet UITextField *messageField;
@property (nonatomic, strong) id<WHOMessageProtocol> delegate;

@end
